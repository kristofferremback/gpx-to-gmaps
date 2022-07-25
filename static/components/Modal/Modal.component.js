/**
 * @param {any} html
 */
const getModal = (html) => {
  const Modal = ({ isOpen, close, title, children }) => {
    return html`
      <dialog class="modal" open=${isOpen}>
        <article>
          <header>
            <a href="#close" aria-label="Close" class="close" onclick=${close}></a>
            ${title}
          </header>
          ${children}
        </article>
      </dialog>
    `
  }
  return Modal
}

export default getModal
