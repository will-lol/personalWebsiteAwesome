class ErrorHandler {
    /**
     * @param {HTMLElement} container
     */
    constructor(container) {
        /**
         * @type {HTMLElement}
         * @private
         */
        this.container = container;
        /** 
         * @type {number}
         * @private
         */
        this.count = 0;
        /**
         * @type {Promise<null>[]}
         * @private
         */
        this.promises = [];
        this.container.classList.add("fixed", "z-50", "right-0", "bottom-0", "p-2", "flex", "gap-2", "flex-col");
    }

    /**
    * @param {string} msg
    * @returns {HTMLDivElement}
    */
    #showError(msg) {
        const elem = document.createElement("div");
        elem.classList.add("border", "bg-warm-900", "pb-2", "pt-1", "px-3", "border-warm-800", "text-warm-100");
        if (msg) {
            elem.innerHTML = "Error: " + msg;
        } else {
            elem.innerHTML = "Something went wrong";
        }
        this.container.appendChild(elem);
        elem.animate([{
            transform: "scale(0.95) translateY(50px)",
            opacity: "0"
        }, {
            transform: "scale(1) translateY(0px)",
            opacity: "1"
        }], {
            duration: 150,
            easing: "ease-out"
        });
        return elem;
    }

    /**
     * @param {string} msg
     */
    enqueue(msg) {
        this.promises.push(new Promise((resolve, _) => {
            if (this.promises.at(-1)) {
                this.promises.at(-1).then(() => {
                    this.#handle(msg, resolve);
                })
            } else {
                this.#handle(msg, resolve);
            }
        }))
    }

    /**
     * @param {string} msg
     * @param {Function} resolve
     */
    #handle(msg, resolve) {
        const el = this.#showError(msg);
        setTimeout(async () => {
            await this.#removeError(el);
            this.promises.pop();
            resolve();
        }, 2000);
    }

    /**
    * @param {HTMLElement} node
    */
    #removeError(node) {
        return new Promise((resolve, _) => {
            node.animate([{
                opacity: "1"
            }, {
                opacity: "0"
            }], {
                duration: 150,
                easing: "ease-in",
                fill: "forwards"
            }).onfinish = () => {
                node.remove();
                resolve();
            };
        })
    }
}

document.addEventListener("htmx:responseError", function(e) {
    handler.enqueue(e.detail.xhr.statusText);
})

document.addEventListener("htmx:sendError", function() {
    handler.enqueue("Couldn't complete request");
})

const handler = new ErrorHandler(document.getElementById("errors"));

export default handler;
