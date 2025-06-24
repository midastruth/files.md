const {test, expect} = require('@playwright/test');
const fs = require('fs').promises;
const path = require('path');
const crypto = require('crypto');

const serverDir = '../storage/-1';

test.beforeEach(async ({page}) => {
    await fs.rm(serverDir, { recursive: true, force: true });
    await fs.mkdir(serverDir, { recursive: true });
    await createFile(saltToken('token'), '-1');
});

async function app(page) {
    await page.addInitScript(() => {
        window.API_HOST = 'http://localhost:8080';
        localStorage.setItem('token', 'token');
    });

    await page.goto('/app.html');

    await page.evaluate(()=> {
        window.getRootDirHandle = async function() {
            const root = await navigator.storage.getDirectory();
            const subdir = await root.getDirectoryHandle('subdir', { create: true });

            const files = [
                { name: 'README.md', content: 'Hello world' },
                { name: 'Notes.md', content: 'Some Text' }
            ];

            for (const file of files) {
                try {
                    await subdir.getFileHandle(file.name);
                } catch (error) {
                    const fileHandle = await subdir.getFileHandle(file.name, { create: true });
                    const writable = await fileHandle.createWritable();
                    await writable.write(file.content);
                    await writable.close();
                }
            }

            return root;
        };
    })
    await page.evaluate(() => {
        init(document.getElementById('editor'));
    });

    await page.waitForSelector('.CodeMirror', {timeout: 10000});
    await page.waitForSelector('#sidebar-tree', {timeout: 5000});
}

test('sync new files from server', async ({ page }) => {
    await createFile('file.md', 'test content');
    await createFile('another.md', '*italic*');

    await app(page);

    await checkFileContent(page, 'subdir/Notes', "# Notes\nSome Text");
    await checkFileContent(page, 'subdir/README', "# README\nHello world");
    // Check that existing files are not removed
    await checkFileContent(page, 'file', "# File\ntest content");
    await checkFileContent(page, 'another', "# Another\n*italic*");
});

// Create file on server.
async function createFile(filepath, content) {
    const p = path.join(serverDir, filepath);
    try {
        await fs.writeFile(p, content, 'utf8');
    } catch (error) {
        console.error('Error creating file:', error);
    }
}

function saltToken(token, salt = "") {
    return crypto.createHash('sha256')
        .update(token + salt)
        .digest('hex');
}

async function checkFileContent(page, filePath, expectedContent) {
    const parts = filePath.split('/');
    const dirs = parts.slice(0, -1);
    const file = parts[parts.length - 1];

    for (const dir of dirs) {
        const isSelected = await page.locator(`#sidebar-tree .tj_description:has-text("${dir}")`).evaluate(el => el.classList.contains('expanded'));
        if (!isSelected) {
            await page.click(`#sidebar-tree .tj_description:has-text("${dir}")`);
            await page.waitForTimeout(100);
        }
    }

    await page.click(`#sidebar-tree .tj_description:has-text("${file}")`);
    await page.waitForTimeout(200);

    const codeMirrorContent = await page.evaluate(() => {
        const cm = document.querySelector('.CodeMirror').CodeMirror;
        return cm.getValue();
    });
    expect(codeMirrorContent).toBe(expectedContent);
}