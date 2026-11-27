import fs from "node:fs";
import path from "node:path";
import ts from "typescript5";

function parseArgs() {
  const args = {
    dts: "dist/types",
    out: "target/sdk.json",
    allowMissing: false,
  };
  for (let i = 2; i < process.argv.length; i++) {
    if (process.argv[i] === "--dts" && process.argv[i + 1]) {
      args.dts = process.argv[++i];
    } else if (process.argv[i] === "--out" && process.argv[i + 1]) {
      args.out = process.argv[++i];
    } else if (process.argv[i] === "--allow-missing") {
      args.allowMissing = true;
    }
  }
  return args;
}

function cleanType(t) {
  if (!t) return "any";
  // Per-file tsc output names types from other files as import("./x.js").Name.
  return t.replace(/import\("[^"]*"\)\./g, "").replace(/\$\d+\b/g, "").trim();
}

function getCommentText(comment) {
  if (!comment) return "";
  let text = "";
  if (typeof comment === "string") {
    text = comment;
  } else if (Array.isArray(comment)) {
    text = comment.map((c) => c.text || "").join("");
  }
  return text.trim().replace(/^-\s+/, "");
}

function parseJSDoc(node) {
  if (!node.jsDoc || node.jsDoc.length === 0) {
    return { summary: "", description: "", tags: [] };
  }
  const last = node.jsDoc[node.jsDoc.length - 1];
  const fullText = getCommentText(last.comment);

  let summary = "";
  let description = "";
  if (fullText) {
    const parts = fullText.split(/\n\s*\n/);
    summary = parts[0].replace(/\n/g, " ").trim();
    if (parts.length > 1) {
      description = parts.slice(1).join("\n\n").trim();
    }
  }

  const tags = [];
  if (last.tags) {
    for (const tag of last.tags) {
      const tagName = tag.tagName?.text || "";
      const name = tag.name?.getText() || tag.name?.text || "";
      const comment = getCommentText(tag.comment);
      tags.push({ tag: tagName, name, comment });
    }
  }

  return { summary, description, tags };
}

// dtsFiles: [{ name, content }], one per .d.ts file.
export function extractDocs(dtsFiles) {
  const classes = new Map();
  const rawTypedefs = [];

  for (const f of dtsFiles) {
    const file = ts.createSourceFile(f.name, f.content, ts.ScriptTarget.ESNext, true);
    for (const s of file.statements) {
      if (ts.isClassDeclaration(s) && s.name) {
        classes.set(s.name.text, s);
      } else if ((ts.isTypeAliasDeclaration(s) || ts.isInterfaceDeclaration(s)) && s.name) {
        rawTypedefs.push([s.name.text, s]);
      }
    }
  }

  const missingDocs = [];

  function checkDoc(name, doc) {
    if (!doc || !doc.trim()) {
      missingDocs.push(`${name}: missing doc`);
    }
  }

  const rootClass = classes.get("I6Sdk");
  if (!rootClass) {
    throw new Error("I6Sdk class not found in declarations");
  }

  // Determine handwritten modules and their call paths
  // I6Sdk is root ("sdk")
  const modulesToProcess = [
    { id: "sdk", name: "I6Sdk", callPath: "sdk", decl: rootClass },
  ];

  for (const member of rootClass.members) {
    if (!ts.isMethodDeclaration(member)) continue;
    const methodName = member.name.getText();
    if (methodName.startsWith("_") || methodName.startsWith("#")) continue;
    if (methodName === "apis" || methodName === "funcs") continue;

    // Check if it's an accessor for a service class: 0 params and returns a class in our file
    if (member.parameters.length === 0 && member.type) {
      const retType = cleanType(member.type.getText());
      const targetClass = classes.get(retType);
      if (targetClass && retType !== "I6Sdk") {
        const accessorDoc = parseJSDoc(member);
        modulesToProcess.push({
          id: methodName.toLowerCase(),
          name: retType,
          callPath: `sdk.${methodName}()`,
          decl: targetClass,
          accessorDoc,
        });
      }
    }
  }

  const modules = [];

  for (const mod of modulesToProcess) {
    const classDoc = parseJSDoc(mod.decl);
    const modSummary = classDoc.summary || mod.accessorDoc?.summary || "";
    const modDescription = classDoc.description || mod.accessorDoc?.description || "";

    const methods = [];
    for (const member of mod.decl.members) {
      if (!ts.isMethodDeclaration(member)) continue;
      const methodName = member.name.getText();
      if (methodName.startsWith("_") || methodName.startsWith("#")) continue;

      // Skip child module accessors on root I6Sdk so they are not listed as methods of sdk
      if (mod.id === "sdk" && member.parameters.length === 0 && member.type) {
        const retType = cleanType(member.type.getText());
        if (classes.has(retType) && retType !== "I6Sdk") continue;
      }

      const mDoc = parseJSDoc(member);
      const callPath = `${mod.callPath}.${methodName}`;
      checkDoc(callPath, mDoc.summary);

      const params = [];
      for (const param of member.parameters) {
        let paramName = param.name.getText();
        const paramType = cleanType(param.type?.getText());
        const optional = param.questionToken != null || param.initializer != null;

        // Match JSDoc @param tag
        let pTag = mDoc.tags.find((t) => t.tag === "param" && t.name === paramName);
        if (!pTag && ts.isObjectBindingPattern(param.name)) {
          // Destructured param: match tag named 'params' or 'opts' or the first @param tag
          pTag = mDoc.tags.find((t) => t.tag === "param");
          if (pTag?.name) {
            paramName = pTag.name;
          }
        }

        const paramDoc = pTag?.comment || "";
        checkDoc(`${callPath} param ${paramName}`, paramDoc);

        params.push({
          name: paramName,
          type: paramType,
          optional,
          doc: paramDoc,
        });
      }

      const retTag = mDoc.tags.find((t) => t.tag === "returns" || t.tag === "return");
      const returnType = cleanType(member.type?.getText()) || "void";
      const returnDoc = retTag?.comment || "";

      methods.push({
        name: methodName,
        callPath,
        summary: mDoc.summary,
        description: mDoc.description,
        params,
        returnType,
        returnDoc,
      });
    }

    modules.push({
      id: mod.id,
      name: mod.name,
      callPath: mod.callPath,
      summary: modSummary,
      description: modDescription,
      methods,
    });
  }

  // Process Typedefs
  const typedefs = {};
  const typedefTexts = {};
  // Texts of the skipped types: what they reference is still documented
  // (on the API pages, or through I6SdkConfig.serviceOptions).
  const skippedTexts = [];
  for (const [rawName, decl] of rawTypedefs) {
    const name = cleanType(rawName);
    // Skip forge API types and internal types
    if (
      name.startsWith("api") ||
      name.startsWith("sdkFunc") ||
      name.startsWith("func") ||
      name.startsWith("Services") ||
      name.startsWith("ServiceOptions") ||
      name.startsWith("dsz")
    ) {
      skippedTexts.push(decl.getText());
      continue;
    }

    typedefTexts[name] = (typedefTexts[name] || "") + decl.getText();
    const tDoc = parseJSDoc(decl);
    const properties = [];

    const targetType = ts.isTypeAliasDeclaration(decl) ? decl.type : decl;
    if (ts.isTypeLiteralNode(targetType) || ts.isInterfaceDeclaration(targetType)) {
      for (const m of targetType.members) {
        if (!ts.isPropertySignature(m)) continue;
        const propName = m.name.getText();
        if (propName.startsWith("_") || propName.startsWith("#")) continue;

        const propDoc = parseJSDoc(m);
        const propType = cleanType(m.type?.getText());
        const optional = m.questionToken != null;

        properties.push({
          name: propName,
          type: propType,
          optional,
          doc: propDoc.summary,
        });
      }
    }

    // Only record if we haven't or if this one has properties
    if (!typedefs[name] || properties.length > 0) {
      typedefs[name] = {
        name,
        summary: tDoc.summary,
        description: tDoc.description,
        properties,
      };
    }
  }

  // Keep only the typedefs reachable from the documented methods, so types
  // used only inside one module (e.g. internal helpers) stay out.
  const reachable = new Set();
  const queue = [
    ...skippedTexts,
    ...modules.flatMap((m) => m.methods.flatMap((x) => [x.returnType, ...x.params.map((p) => p.type)])),
  ];
  while (queue.length > 0) {
    for (const id of queue.pop().match(/[A-Za-z_]\w*/g) || []) {
      if (typedefs[id] && !reachable.has(id)) {
        reachable.add(id);
        queue.push(typedefTexts[id]);
      }
    }
  }
  for (const name of Object.keys(typedefs)) {
    if (!reachable.has(name)) {
      delete typedefs[name];
      continue;
    }
    for (const prop of typedefs[name].properties) {
      checkDoc(`${name}.${prop.name}`, prop.doc);
    }
  }

  return {
    data: {
      modules,
      typedefs,
    },
    missingDocs,
  };
}

// A single .d.ts file, or every .d.ts under a directory.
function readDtsFiles(p) {
  if (!fs.statSync(p).isDirectory()) {
    return [{ name: p, content: fs.readFileSync(p, "utf8") }];
  }
  return fs
    .readdirSync(p, { recursive: true })
    .filter((f) => f.endsWith(".d.ts"))
    .sort()
    .map((f) => ({ name: f, content: fs.readFileSync(path.join(p, f), "utf8") }));
}

function main() {
  const args = parseArgs();
  if (!fs.existsSync(args.dts)) {
    console.error(`docgen: file not found: ${args.dts}`);
    process.exit(1);
  }

  const { data, missingDocs } = extractDocs(readDtsFiles(args.dts));

  if (missingDocs.length > 0) {
    for (const item of missingDocs) {
      console.error(item);
    }
    if (!args.allowMissing) {
      process.exit(1);
    }
  }

  fs.mkdirSync(path.dirname(args.out), { recursive: true });
  fs.writeFileSync(args.out, JSON.stringify(data, null, 2) + "\n", "utf8");
}

if (process.argv[1] && import.meta.url.endsWith(path.basename(process.argv[1]))) {
  main();
}
