// search.js — Real-time sidebar navigation filter

(function () {
  function initSearch() {
    const input = document.getElementById("search-input");
    const sidebar = document.getElementById("sidebar");
    if (!input || !sidebar) return;

    input.addEventListener("input", () => {
      const q = input.value.trim().toLowerCase();
      const sections = sidebar.querySelectorAll(".nav-section");

      if (!q) {
        // Reset all visibility
        sidebar.querySelectorAll(".nav-item, .nav-subitem, .nav-section").forEach((el) => {
          el.style.display = "";
        });
        return;
      }

      sections.forEach((section) => {
        let sectionHasMatch = false;

        const items = section.querySelectorAll(".nav-item");
        items.forEach((item) => {
          const itemText = item.querySelector("a")?.textContent.toLowerCase() || "";
          const subitems = item.querySelectorAll(".nav-subitem");
          let subHasMatch = false;

          subitems.forEach((sub) => {
            const subText = sub.textContent.toLowerCase();
            const match = subText.includes(q);
            sub.style.display = match ? "" : "none";
            if (match) subHasMatch = true;
          });

          const itemMatches = itemText.includes(q);
          const shouldShow = itemMatches || subHasMatch;
          item.style.display = shouldShow ? "" : "none";
          if (shouldShow) sectionHasMatch = true;
        });

        // Also check direct subitems in sections without nav-item (like APIs)
        const directSubitems = section.querySelectorAll(".nav-list > .nav-subitem");
        directSubitems.forEach((sub) => {
          const text = sub.textContent.toLowerCase();
          const match = text.includes(q);
          sub.style.display = match ? "" : "none";
          if (match) sectionHasMatch = true;
        });

        section.style.display = sectionHasMatch ? "" : "none";
      });
    });

    // Press Escape to clear search
    input.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        input.value = "";
        input.dispatchEvent(new Event("input"));
        input.blur();
      }
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initSearch);
  } else {
    initSearch();
  }
})();
