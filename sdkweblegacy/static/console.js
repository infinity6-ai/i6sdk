// Mirrors console output and uncaught errors into #console-log on the page.
(function () {
    const log = document.getElementById("console-log")
    const levels = ["log", "info", "warn", "error"]

    function fmt(arg) {
        if (arg instanceof Error) return arg.stack || String(arg)
        if (typeof arg === "string") return arg
        try {
            return JSON.stringify(arg, null, 2) ?? String(arg)
        } catch {
            return String(arg)
        }
    }

    function add(level, args) {
        const line = document.createElement("div")
        line.className = "line " + level
        const time = document.createElement("span")
        time.className = "time"
        time.textContent = new Date().toLocaleTimeString()
        line.append(time, args.map(fmt).join(" "))
        log.append(line)
        log.scrollTop = log.scrollHeight
    }

    for (const level of levels) {
        const orig = console[level].bind(console)
        console[level] = (...args) => {
            orig(...args)
            add(level, args)
        }
    }

    window.addEventListener("error", (e) => add("error", [e.error ?? e.message]))
    window.addEventListener("unhandledrejection", (e) => add("error", ["Unhandled rejection:", e.reason]))

    document.getElementById("console-clear").addEventListener("click", () => log.replaceChildren())
})()
