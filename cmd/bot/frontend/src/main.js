import './style.css';

import {Send} from '../wailsjs/go/main/App';
import {NewUpdate} from '../wailsjs/go/main/App';
import {NewCmd} from '../wailsjs/go/main/App';

window.send = Send
window.newUpdate = NewUpdate
window.newCmd = NewCmd