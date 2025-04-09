/**
 * Guestbook Application
 * A modular JavaScript application for managing guestbook messages
 */

// API Service - Handles all communication with the backend
const APIService = {
    baseUrl: '/api/v1/messages',

    async fetchMessages() {
        const response = await fetch(this.baseUrl);
        if (!response.ok) throw new Error('Failed to fetch messages');
        return response.json();
    },

    async addMessage(author, content) {
        return fetch(this.baseUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ author, content })
        });
    },

    async deleteMessage(id) {
        return fetch(`${this.baseUrl}/${id}`, { method: 'DELETE' });
    }
};

// Form Validator - Handles form validation logic
const FormValidator = {
    validateField(value, fieldName) {
        if (!value || value.trim() === '') {
            return {
                valid: false,
                message: `${fieldName} cannot be empty`
            };
        }
        return { valid: true };
    },

    validateForm(nameValue, messageValue) {
        const nameValidation = this.validateField(nameValue, 'Name');
        const messageValidation = this.validateField(messageValue, 'Message');

        return {
            valid: nameValidation.valid && messageValidation.valid,
            nameError: nameValidation.valid ? null : nameValidation.message,
            messageError: messageValidation.valid ? null : messageValidation.message
        };
    }
};

// UI Manager - Handles all UI-related operations
const UIManager = {
    elements: {
        form: document.getElementById('messageForm'),
        nameInput: document.getElementById('name'),
        messageInput: document.getElementById('message'),
        messagesContainer: document.getElementById('messages')
    },

    toggleError(input, message) {
        if (message) {
            input.placeholder = message;
            input.classList.add('error-placeholder');
            input.value = '';
        } else {
            input.classList.remove('error-placeholder');
        }
    },

    clearErrors() {
        this.toggleError(this.elements.nameInput, null);
        this.toggleError(this.elements.messageInput, null);
    },

    displayMessages(messages) {
        const container = this.elements.messagesContainer;
        container.innerHTML = '';

        if (!messages || messages.length === 0) {
            container.innerHTML = '<p>No messages to display.</p>';
            return;
        }

        const fragment = document.createDocumentFragment();
        messages.forEach(msg => {
            const div = document.createElement('div');
            div.className = 'message';
            div.innerHTML = `
                <span class="delete-btn" data-id="${msg.id}">&times;</span>
                <strong>${msg.author || 'Anonymous'}</strong>
                <p>${msg.content || 'No content provided'}</p>
            `;
            fragment.appendChild(div);
        });
        container.appendChild(fragment);
    },

    clearMessageInput() {
        this.elements.messageInput.value = '';
    },

    showError(message) {
        console.error(message);
        this.elements.messagesContainer.innerHTML =
            `<p class="error-message">Error: ${message}. Please try again later.</p>`;
    }
};

// Event Handler - Manages all event-related logic
const EventHandler = {
    init() {
        // Bind methods to preserve context
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleMessageKeyPress = this.handleMessageKeyPress.bind(this);
        this.handleInput = this.handleInput.bind(this);
        this.handleMessageClick = this.handleMessageClick.bind(this);

        // Form submission
        UIManager.elements.form.addEventListener('submit', this.handleSubmit);

        // Input events
        UIManager.elements.messageInput.addEventListener('keypress', this.handleMessageKeyPress);
        UIManager.elements.messageInput.addEventListener('input', this.handleInput);
        UIManager.elements.nameInput.addEventListener('input', this.handleInput);

        // Focus events
        UIManager.elements.messageInput.addEventListener('focus',
            () => UIManager.toggleError(UIManager.elements.messageInput, null));
        UIManager.elements.nameInput.addEventListener('focus',
            () => UIManager.toggleError(UIManager.elements.nameInput, null));

        // Message container for delete buttons
        UIManager.elements.messagesContainer.addEventListener('click', this.handleMessageClick);
    },

    async handleSubmit(e) {
        e.preventDefault();
        const author = UIManager.elements.nameInput.value.trim();
        const content = UIManager.elements.messageInput.value.trim();

        const validation = FormValidator.validateForm(author, content);
        if (!validation.valid) {
            if (validation.nameError) UIManager.toggleError(UIManager.elements.nameInput, validation.nameError);
            if (validation.messageError) UIManager.toggleError(UIManager.elements.messageInput, validation.messageError);
            return;
        }

        try {
            const response = await APIService.addMessage(author, content);
            if (response.ok) {
                UIManager.clearMessageInput();
                GuestbookController.loadMessages();
            } else {
                UIManager.showError('Failed to add message');
            }
        } catch (error) {
            UIManager.showError(error.message);
        }
    },

    handleMessageKeyPress(e) {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            UIManager.elements.form.dispatchEvent(new Event('submit'));
        }
    },

    handleInput(e) {
        if (e.target.value.trim() !== '') {
            e.target.classList.remove('error-placeholder');
        }
    },

    async handleMessageClick(e) {
        if (e.target.classList.contains('delete-btn')) {
            const messageId = e.target.dataset.id;
            try {
                const response = await APIService.deleteMessage(messageId);
                if (response.ok) {
                    GuestbookController.loadMessages();
                } else {
                    UIManager.showError('Failed to delete message');
                }
            } catch (error) {
                UIManager.showError(error.message);
            }
        }
    }
};

// Main Application Controller
const GuestbookController = {
    async loadMessages() {
        try {
            const data = await APIService.fetchMessages();

            if (!data.messages || !Array.isArray(data.messages)) {
                throw new Error('Unexpected response format');
            }

            UIManager.displayMessages(data.messages);
        } catch (error) {
            UIManager.showError(error.message);
        }
    },

    init() {
        EventHandler.init();
        this.loadMessages();
    }
};

// Initialize the application when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    GuestbookController.init();
});
