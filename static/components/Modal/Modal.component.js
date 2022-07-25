import html from '../../lib/html.js'

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

export default Modal
