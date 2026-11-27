// bar.js — SDK config bar and account authentication state management

const STORAGE_KEY = "i6_docs_sdk_config";

// Docs root (where callback.html lives), taken from this script's own URL so
// it is right on nested pages like /apis/x.html too.
const DOCS_ROOT_URL = new URL(".", (document.currentScript && document.currentScript.src) || window.location.href).href;

// The docs run with the SDK defaults (app, servUrls); only the sign-in
// callback points at this site.
function getDefaultSdkConfig() {
  return {
    userExchangeCallbackUrl: new URL("callback.html", DOCS_ROOT_URL).href,
  };
}

function getSavedSdkConfig() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      return { ...getDefaultSdkConfig(), ...parsed };
    }
  } catch (e) {
    console.warn("Failed to load saved SDK config:", e);
  }
  return getDefaultSdkConfig();
}

function saveSdkConfig(config) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(config));
  } catch (e) {
    console.error("Failed to save SDK config:", e);
  }
}

function resetSdkConfig() {
  try {
    localStorage.removeItem(STORAGE_KEY);
  } catch (e) {
    console.error("Failed to clear saved SDK config:", e);
  }
}

function formatConfigForPreview(config) {
  const lines = [
    "const sdk = new I6Sdk({",
    `  app: ${JSON.stringify(config.app || "default")},`,
  ];
  if (config.servUrls) {
    lines.push("  servUrls: {");
    for (const [k, v] of Object.entries(config.servUrls)) {
      lines.push(`    ${k}: ${JSON.stringify(v)},`);
    }
    lines.push("  },");
  }
  if (config.userExchangeCallbackUrl) {
    lines.push(`  userExchangeCallbackUrl: ${JSON.stringify(config.userExchangeCallbackUrl)},`);
  }
  lines.push("})");
  return lines.join("\n");
}

function updateTopBarUI(config) {
  const preview = document.getElementById("sdk-config-preview");
  if (preview) {
    // Show what the SDK really runs with (its defaults merged with ours).
    const effective = typeof I6Sdk !== "undefined" ? new I6Sdk(config).config() : config;
    preview.textContent = formatConfigForPreview(effective);
  }

  const resetBtn = document.getElementById("btn-reset-config");
  if (resetBtn) {
    const isCustom = localStorage.getItem(STORAGE_KEY) !== null;
    resetBtn.style.display = isCustom ? "inline-block" : "none";
  }
}

function openConfigEditorModal(currentConfig, onSave) {
  const overlay = document.createElement("div");
  overlay.className = "modal-overlay";

  const content = document.createElement("div");
  content.className = "modal-content";

  const header = document.createElement("div");
  header.className = "modal-header";
  const title = document.createElement("h3");
  title.className = "modal-title";
  title.textContent = "Edit SDK Configuration";
  header.appendChild(title);

  const textarea = document.createElement("textarea");
  textarea.className = "modal-textarea";
  textarea.value = JSON.stringify(currentConfig, null, 2);

  const actions = document.createElement("div");
  actions.className = "modal-actions";

  const cancelBtn = document.createElement("button");
  cancelBtn.className = "btn btn-secondary";
  cancelBtn.type = "button";
  cancelBtn.textContent = "Cancel";
  cancelBtn.onclick = () => document.body.removeChild(overlay);

  const saveBtn = document.createElement("button");
  saveBtn.className = "btn btn-primary";
  saveBtn.type = "button";
  saveBtn.textContent = "Save";
  saveBtn.onclick = () => {
    try {
      const parsed = JSON.parse(textarea.value);
      onSave(parsed);
      document.body.removeChild(overlay);
    } catch (err) {
      alert("Invalid JSON: " + err.message);
    }
  };

  actions.appendChild(cancelBtn);
  actions.appendChild(saveBtn);

  content.appendChild(header);
  content.appendChild(textarea);
  content.appendChild(actions);
  overlay.appendChild(content);

  document.body.appendChild(overlay);
  textarea.focus();
}

async function updateAccountUI(sdk) {
  const label = document.getElementById("account-label");
  const authBtn = document.getElementById("btn-auth");
  if (!label || !authBtn) return;

  try {
    const token = sdk.auth().getUserToken();
    if (token && token.user_access_token) {
      let userDesc = "Signed in";
      try {
        const u = await sdk.auth().user();
        if (u) userDesc = u.email_redacted || u.email || u.display || u.user_id || userDesc;
      } catch (e) {
        console.warn("Could not load signed-in user:", e);
      }
      label.textContent = userDesc;
      label.title = "Authenticated";
      authBtn.textContent = "Sign out";
      authBtn.onclick = async () => {
        try {
          // Clears the tokens and navigates back to this page.
          await sdk.auth().logout(window.location.href);
        } catch (err) {
          console.warn("Logout call failed, clearing local tokens:", err);
          sdk.auth().clearTokens();
          window.location.reload();
        }
      };
      return;
    }
  } catch (err) {
    console.warn("Auth check failed:", err);
  }

  label.textContent = "Not signed in";
  authBtn.textContent = "Sign in";
  authBtn.onclick = async () => {
    try {
      await sdk.auth().requireUser();
    } catch (err) {
      alert("Sign in redirect failed: " + (err.message || err));
    }
  };
}

function initBar() {
  const config = getSavedSdkConfig();
  updateTopBarUI(config);

  if (typeof I6Sdk !== "undefined" && !window.sdk) {
    window.sdk = new I6Sdk(config);
    updateAccountUI(window.sdk);
  }

  const editBtn = document.getElementById("btn-edit-config");
  if (editBtn) {
    editBtn.onclick = () => {
      openConfigEditorModal(getSavedSdkConfig(), (newConfig) => {
        saveSdkConfig(newConfig);
        updateTopBarUI(newConfig);
        if (typeof I6Sdk !== "undefined") {
          window.sdk = new I6Sdk(newConfig);
          updateAccountUI(window.sdk);
        }
      });
    };
  }

  const resetBtn = document.getElementById("btn-reset-config");
  if (resetBtn) {
    resetBtn.onclick = () => {
      resetSdkConfig();
      const def = getDefaultSdkConfig();
      updateTopBarUI(def);
      if (typeof I6Sdk !== "undefined") {
        window.sdk = new I6Sdk(def);
        updateAccountUI(window.sdk);
      }
    };
  }
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", initBar);
} else {
  initBar();
}

window.getSavedSdkConfig = getSavedSdkConfig;
