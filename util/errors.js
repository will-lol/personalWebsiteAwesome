class ErrorHandler {
    /**
     * @param {HTMLElement} container
     */
    constructor(container) {
        this.container = container;
        this.container.classList.add("fixed", "z-50", "right-0", "bottom-0", "p-2", "flex", "gap-2", "flex-col");
    }

    /**
    * @param {string} msg
    */
    showError(msg) {
        const elem = document.createElement("div");
        console.log("hi");
        elem.classList.add("border", "bg-warm-900", "pb-2", "pt-1", "px-3", "border-warm-800", "text-warm-100");
        if (msg) {
            elem.innerHTML = "Error: " + msg;
        } else {
            elem.innerHTML = "Something went wrong";
        }
        this.container.appendChild(elem);
        setTimeout(() => this.removeError(elem), 2000)
    }

    /**
    * @param {HTMLElement} node
    */
    removeError(node) {
        node.remove();
    }
}

document.addEventListener("htmx:responseError", function(e) {
    handler.showError(e.detail.xhr.statusText);
})

document.addEventListener("htmx:sendError", function() {
    handler.showError("Couldn't complete request");
})

const handler = new ErrorHandler(document.getElementById("errors"));

export default handler;
