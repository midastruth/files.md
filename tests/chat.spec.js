const {test, expect} = require('@playwright/test');

test.beforeEach(async ({page}) => {
    await page.goto('/app.html');

    await page.waitForSelector('.CodeMirror', {timeout: 10000});
    await page.waitForSelector('#sidebar-tree', {timeout: 5000});
});

test('toggle chat mode', async ({ page }) => {
    await page.evaluate(() => {
        window.getRootDirHandle = async function() {
            const opfsRoot = await navigator.storage.getDirectory();
            const testDir = await opfsRoot.getDirectoryHandle('test-files', { create: true });

            const fileHandle = await testDir.getFileHandle('Notes.md', { create: true });
            const writable = await fileHandle.createWritable();
            await writable.write('Some notes content');
            await writable.close();

            return testDir;
        };
    });

    await page.evaluate(() => {
        init(document.getElementById("editor"));
    });

    await page.waitForTimeout(500);

    // Initially should be in editor mode
    expect(await page.isVisible('#content')).toBe(true);
    expect(await page.isVisible('#sidebar-container')).toBe(true);

    // Chat container should not exist yet or be hidden
    const chatContainerExists = await page.locator('#chat-container').count() > 0;
    if (chatContainerExists) {
        expect(await page.isVisible('#chat-container')).toBe(false);
    }

    // Verify isChat is false initially
    const initialChatMode = await page.evaluate(() => isChat);
    expect(initialChatMode).toBe(false);

    // Toggle to chat mode with Meta+Enter
    await page.keyboard.press('Meta+Enter');
    await page.waitForTimeout(200);

    // Verify chat mode is activated
    const chatModeAfterToggle = await page.evaluate(() => isChat);
    expect(chatModeAfterToggle).toBe(true);

    // Check that UI has switched
    expect(await page.isVisible('#sidebar-container')).toBe(false);
    expect(await page.isVisible('#content')).toBe(false);

    // Chat container should be visible now
    const chatContainerVisible = await page.isVisible('#chat-container');
    expect(chatContainerVisible).toBe(true);

    // Toggle back to editor mode with Meta+Enter
    await page.keyboard.press('Meta+Enter');
    await page.waitForTimeout(200);

    // Verify back in editor mode
    const editorModeAfterToggle = await page.evaluate(() => isChat);
    expect(editorModeAfterToggle).toBe(false);

    // UI should be back to editor
    expect(await page.isVisible('#content')).toBe(true);
    expect(await page.isVisible('#sidebar-container')).toBe(true);
    expect(await page.isVisible('#chat-container')).toBe(false);

    // Verify editor has focus after returning from chat
    const editorFocused = await page.evaluate(() => {
        const cm = document.querySelector('.CodeMirror').CodeMirror;
        return cm.hasFocus();
    });
    expect(editorFocused).toBe(true);
});

test('send chat message and it is in saved', async ({ page }) => {
    await page.evaluate(() => {
        window.getRootDirHandle = async function() {
            const opfsRoot = await navigator.storage.getDirectory();

            return opfsRoot;
        };
    });

    await page.evaluate(() => {
        init(document.getElementById("editor"));
    });

    await page.waitForTimeout(500);

    // Switch to chat mode
    await page.keyboard.press('Meta+Enter');
    await page.waitForTimeout(200);

    // Verify we're in chat mode
    expect(await page.isVisible('#chat-container')).toBe(true);
    expect(await page.isVisible('#content')).toBe(false);

    // Type a message in the chat input
    const chatInput = page.locator('#input-field');
    await chatInput.click();
    await chatInput.fill('My thoughts');

    await page.keyboard.press('Enter');
    await page.waitForTimeout(500);

    // Switch back to editor mode
    await page.keyboard.press('Meta+Enter');
    await page.waitForTimeout(200);

    expect(await page.isVisible('#content')).toBe(true);
    expect(await page.isVisible('#chat-container')).toBe(false);

    await page.waitForTimeout(200);

    await page.click('#sidebar >> text=Saved');
    await page.waitForTimeout(200);

    // Get the updated content of the today file
    const updatedContent = await page.evaluate(() => {
        const cm = document.querySelector('.CodeMirror').CodeMirror;
        return cm.getValue();
    });


    expect(updatedContent).toContain('# Saved\nMy thoughts');
});

// test('send chat message and verify task in today folder', async ({ page }) => {
//     await page.evaluate(() => {
//         window.getRootDirHandle = async function() {
//             const opfsRoot = await navigator.storage.getDirectory();
//             const testDir = await opfsRoot.getDirectoryHandle('test-files', { create: true });
//
//             // Create today directory with initial content
//             const todayDir = await testDir.getDirectoryHandle('today', { create: true });
//             const todayFile = await todayDir.getFileHandle('Task.md', { create: true });
//             const writable = await todayFile.createWritable();
//             await writable.write('');
//             await writable.close();
//
//             return testDir;
//         };
//     });
//
//     await page.evaluate(() => {
//         init(document.getElementById("editor"));
//     });
//
//     await page.waitForTimeout(500);
//
//     // Check initial today folder content
//     await page.click('#sidebar >> text=today');
//     await page.waitForTimeout(100);
//     await page.click('#sidebar >> text=Task');
//     await page.waitForTimeout(200);
//
//     const initialContent = await page.evaluate(() => {
//         const cm = document.querySelector('.CodeMirror').CodeMirror;
//         return cm.getValue();
//     });
//     expect(initialContent).toContain('# Task');
//
//     // Switch to chat mode
//     await page.keyboard.press('Meta+Enter');
//     await page.waitForTimeout(200);
//
//     // Verify we're in chat mode
//     expect(await page.isVisible('#chat-container')).toBe(true);
//     expect(await page.isVisible('#content')).toBe(false);
//
//     // Type a message in the chat input
//     const chatInput = page.locator('#input-field');
//     await chatInput.click();
//     await chatInput.fill('Task2');
//
//     await page.keyboard.press('Enter');
//     await page.waitForTimeout(500);
//
//     // Switch back to editor mode
//     await page.keyboard.press('Meta+Enter');
//     await page.waitForTimeout(200);
//
//     expect(await page.isVisible('#content')).toBe(true);
//     expect(await page.isVisible('#chat-container')).toBe(false);
//
//     await page.waitForTimeout(200);
//
//     await page.click('#sidebar >> text=Task2');
//     await page.waitForTimeout(200);
//
//     // Get the updated content of the today file
//     const updatedContent = await page.evaluate(() => {
//         const cm = document.querySelector('.CodeMirror').CodeMirror;
//         return cm.getValue();
//     });
//
//
//     expect(updatedContent).toContain('# Task2');
// });