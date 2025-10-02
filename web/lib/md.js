// Various string functions, ported from Golang bot.

// Add content at the beginning of the file, prepending current's day header
async function addToFile(path, content) {
    let existingContent = '';

    existingContent = read(path);

    const now = new Date();
    const monthNames = [
        'January', 'February', 'March', 'April', 'May', 'June',
        'July', 'August', 'September', 'October', 'November', 'December'
    ];
    const dayNames = [
        'Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'
    ];

    const header = `#### ${now.getDate()} ${monthNames[now.getMonth()]} ${now.getFullYear()}, ${dayNames[now.getDay()]}`;
    const newContent = insertTextAfterHeader(existingContent, header, content);

    await write(path, newContent);
}

function insertTextAfterHeader(existingContent, header, newContent) {
    if (!existingContent.includes(header)) {
        if (existingContent === "") {
            return `${header}\n${newContent}`;
        } else {
            return `${header}\n${newContent}\n\n${existingContent}`;
        }
    }

    const lines = existingContent.split("\n");
    let headerIndex = -1;

    // Find the header line
    for (let i = 0; i < lines.length; i++) {
        if (lines[i] === header) {
            headerIndex = i;
            break;
        }
    }

    if (headerIndex === -1) {
        return `${header}\n${newContent}\n\n${existingContent}`;
    }

    // Find where to insert (after the last line belonging to this header)
    let insertIndex = headerIndex + 1;

    // Look for the next header or end of content
    for (let i = headerIndex + 1; i < lines.length; i++) {
        if (lines[i].startsWith("###")) {
            insertIndex = i;
            break;
        }
        // If we encounter an empty line, insert before it
        if (lines[i].trim() === "") {
            insertIndex = i;
            break;
        }
        insertIndex = i + 1;
    }

    // Insert the new content
    const newLines = [];
    newLines.push(...lines.slice(0, insertIndex));
    newLines.push(newContent);

    // Add empty line after new content if there's content following and it's not empty
    if (insertIndex < lines.length && lines[insertIndex].trim() !== "") {
        newLines.push("");
    }

    newLines.push(...lines.slice(insertIndex));

    return newLines.join("\n");
}

async function addToJournal(record) {
    record = record.trim();
    const journalFilename = todayJournalFilename();
    const journalPath = `journal/${journalFilename}`;

    let md = '';
    try {
        md = read(journalPath);
        md = normNewLines(md);
        md = md.trim();
        if (md.length !== 0) {
            md += '\n';
        }
    } catch (err) {
        // File doesn't exist, will be created
        md = '';
    }

    const header = todayJournalHeader();
    if (!md.includes(header)) {
        md += header + '\n';
    }

    const now = new Date();
    const timestamp = `\`${now.toLocaleTimeString('en-US', {
        hour12: false,
        hour: '2-digit',
        minute: '2-digit'
    })}\``;

    if (hasImage(record)) {
        // If there's an image - place timestamp under the image
        const imgMatch = record.match(IMG_PATTERN);
        if (imgMatch) {
            const imgLink = imgMatch[0];
            record = record.replace(imgLink, '').trim();
            record = `${imgLink}\n${timestamp} ${record.trim()}\n`;
        }
    } else {
        record = `${timestamp} ${record}\n`;
    }

    md += record;

    await write(journalPath, md);
}

function todayJournalFilename(timezone) {
    const now = new Date();
    const year = now.toLocaleDateString('en-US', { year: 'numeric'});
    const month = now.toLocaleDateString('en-US', { month: '2-digit', timeZone: timezone });
    return `${year}-${month}.md`;
}

function todayJournalHeader(timezone) {
    const now = new Date();
    const monthNames = [
        'January', 'February', 'March', 'April', 'May', 'June',
        'July', 'August', 'September', 'October', 'November', 'December'
    ];
    const dayNames = [
        'Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'
    ];

    const day = parseInt(now.toLocaleDateString('en-US', { day: 'numeric', timeZone: timezone }));
    const monthIndex = parseInt(now.toLocaleDateString('en-US', { month: 'numeric', timeZone: timezone })) - 1;
    const year = parseInt(now.toLocaleDateString('en-US', { year: 'numeric', timeZone: timezone }));
    const dayIndex = new Date(now.toLocaleDateString('en-US', { timeZone: timezone })).getDay();

    return `#### ${day} ${monthNames[monthIndex]} ${year}, ${dayNames[dayIndex]}`;
}

function normNewLines(text) {
    return text.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
}

function hasImage(text) {
    return IMG_PATTERN.test(text);
}

// Define the image pattern constant
const IMG_PATTERN = /!\[.*?\]\(.*?\)/;

