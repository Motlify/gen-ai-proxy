// src/templates/ui/js/modal.js

let genericModal, closeModalButton, modalTitle, modalContent;

export const initModal = (modalElement, closeButton, titleElement, contentElement) => {
    genericModal = modalElement;
    closeModalButton = closeButton;
    modalTitle = titleElement;
    modalContent = contentElement;

    if (closeModalButton) {
        closeModalButton.addEventListener('click', closeModal);
        genericModal.addEventListener('click', (e) => {
            if (e.target === genericModal) {
                closeModal();
            }
        });
    }
};

export const openModal = (title, formHtml) => {
    modalTitle.textContent = title;
    modalContent.innerHTML = formHtml;
    genericModal.classList.remove('hidden');
};

export const closeModal = () => {
    genericModal.classList.add('hidden');
    modalTitle.textContent = '';
    modalContent.innerHTML = '';
};
