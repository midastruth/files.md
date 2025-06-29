class SearchModal {
    constructor() {
        this.mode = 'default';
        this.messageIndex = null;
        this.focusedIndex = 0;
        this.init();
    }

    init() {
        document.getElementById('search').addEventListener('keydown', (event) => {
            const resultsList = document.getElementById('search-results').querySelectorAll('li');

            if (event.key === 'Enter') {
                event.preventDefault();
                this.handleEnterKey();
            }

            if (event.key === 'ArrowDown') {
                event.preventDefault();
                this.focusedIndex = (this.focusedIndex + 1) % resultsList.length;
                this.updateFocusedItem();
            } else if (event.key === 'ArrowUp') {
                event.preventDefault();
                this.focusedIndex = (this.focusedIndex - 1 + resultsList.length) % resultsList.length;
                this.updateFocusedItem();
            }
        });
    }

    open(text = '', mode = 'default', messageIndex = null) {
        this.mode = mode;
        this.messageIndex = messageIndex;

        document.getElementById('search').style.display = 'block';
        const inputField = document.getElementById('search-input');
        inputField.value = text;
        inputField.focus();

        this.focusedIndex = 0;
        const goToFileResults = document.getElementById('search-results');
        goToFileResults.innerHTML = '';

        if (text === '') {
            loadRecentFiles();
        } else {
            search();
        }
    }

    close() {
        document.getElementById('search').style.display = 'none';
        this.mode = 'default';
        this.messageIndex = null;
    }

    showResults(results) {
        const list = document.getElementById('search-results');
        list.innerHTML = ''; // Clear previous results

        results.forEach(({dir, filename}, index) => {
            if (filename === CONFIG_FILENAME) {
                return;
            }

            const listItem = document.createElement('li');
            let title = filename.replace(/\.md$/, '')
            if (dir !== '') {
                listItem.textContent = `${dir}/${title}`;
            } else {
                listItem.textContent = title;
            }
            listItem.setAttribute('data-path', `${dir}/${filename}`);
            listItem.setAttribute('data-index', index);

            listItem.onclick = () => this.handleItemClick(dir, filename);

            listItem.onmouseenter = () => {
                document.querySelectorAll('#search-results li').forEach(li => li.classList.remove('focused'));
                listItem.classList.add('focused');
                this.focusedIndex = index;
            };
            list.appendChild(listItem);
        });

        this.focusedIndex = 0;
        this.updateFocusedItem();
    }

    handleItemClick(dir, filename) {
        if (this.mode === 'move-file') {
            let cmd = {
                n: "mf",
                t: "cmd",
                p: [filename, this.messageIndex.toString()]
            }
            replyCmd(JSON.stringify(cmd));
            this.close();
        } else {
            // Default behavior
            openEditor(!isChat);
            openFile(dir, filename);
            this.close();
        }
    }

    handleEnterKey() {
        const resultsList = document.getElementById('search-results').querySelectorAll('li');
        if (resultsList[this.focusedIndex]) {
            const [dir, filename] = resultsList[this.focusedIndex].getAttribute('data-path').split('/');
            this.handleItemClick(dir, filename);
        }
    }

    updateFocusedItem() {
        const resultsList = document.getElementById('search-results').querySelectorAll('li');
        document.querySelectorAll('#search-results li').forEach(li => li.classList.remove('focused'));
        resultsList.forEach((item, index) => {
            if (index === this.focusedIndex) {
                item.classList.add('focused');
                item.scrollIntoView({block: 'nearest'});
            } else {
                item.classList.remove('focused');
            }
        });
    }
}

const searchModal = new SearchModal();