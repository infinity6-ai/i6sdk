import fs from "node:fs";
import path from "node:path";
import ts from "typescript5";

function parseArgs() {
  const args = {
    dts: "dist/bundle/i6sdk-web.d.ts",
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
  return t.replace(/\$\d+\b/g, "").trim();
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

export function extractDocs(dtsContent) {
  const file = ts.createSourceFile("bundle.d.ts", dtsContent, ts.ScriptTarget.ESNext, true);

  const classes = new Map();
  const rawTypedefs = new Map();

  for (const s of file.statements) {
    if (ts.isClassDeclaration(s) && s.name) {
      classes.set(s.name.text, s);
    } else if ((ts.isTypeAliasDeclaration(s) || ts.isInterfaceDeclaration(s)) && s.name) {
      rawTypedefs.set(s.name.text, s);
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
    const methodName = member.name.getText(file);
    if (methodName.startsWith("_") || methodName.startsWith("#")) continue;
    if (methodName === "apis" || methodName === "funcs") continue;

    // Check if it's an accessor for a service class: 0 params and returns a class in our file
    if (member.parameters.length === 0 && member.type) {
      const retType = cleanType(member.type.getText(file));
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
      const methodName = member.name.getText(file);
      if (methodName.startsWith("_") || methodName.startsWith("#")) continue;

      // Skip child module accessors on root I6Sdk so they are not listed as methods of sdk
      if (mod.id === "sdk" && member.parameters.length === 0 && member.type) {
        const retType = cleanType(member.type.getText(file));
        if (classes.has(retType) && retType !== "I6Sdk") continue;
      }

      const mDoc = parseJSDoc(member);
      const callPath = `${mod.callPath}.${methodName}`;
      checkDoc(callPath, mDoc.summary);

      const params = [];
      for (const param of member.parameters) {
        let paramName = param.name.getText(file);
        const paramType = cleanType(param.type?.getText(file));
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
      const returnType = cleanType(member.type?.getText(file)) || "void";
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
  for (const [rawName, decl] of rawTypedefs.entries()) {
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
      continue;
    }

    const tDoc = parseJSDoc(decl);
    const properties = [];

    const targetType = ts.isTypeAliasDeclaration(decl) ? decl.type : decl;
    if (ts.isTypeLiteralNode(targetType) || ts.isInterfaceDeclaration(targetType)) {
      for (const m of targetType.members) {
        if (!ts.isPropertySignature(m)) continue;
        const propName = m.name.getText(file);
        if (propName.startsWith("_") || propName.startsWith("#")) continue;

        const propDoc = parseJSDoc(m);
        const propType = cleanType(m.type?.getText(file));
        const optional = m.questionToken != null;

        checkDoc(`${name}.${propName}`, propDoc.summary);

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

  return {
    data: {
      modules,
      typedefs,
    },
    missingDocs,
  };
}

function main() {
  const args = parseArgs();
  if (!fs.existsSync(args.dts)) {
    console.error(`docgen: file not found: ${args.dts}`);
    process.exit(1);
  }

  const dtsContent = fs.readFileSync(args.dts, "utf8");
  const { data, missingDocs } = extractDocs(dtsContent);

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
