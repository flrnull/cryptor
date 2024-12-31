import './style.css';

document.querySelector('#app').innerHTML = `
  <div style="text-align: center; margin: 20px;">
    <h1>Cryptor</h1>
    <p>Use this app to encrypt or decrypt your files.</p>
    <button id="openFile">Open File</button>
    <textarea id="fileContent" rows="10" cols="80" placeholder="File content will appear here..."></textarea>
    <div style="margin-top: 20px;">
      <input type="password" id="password" placeholder="Enter password" />
      <button id="encryptFile">Encrypt</button>
      <button id="decryptFile">Decrypt</button>
    </div>
  </div>
`;

// Функция для взаимодействия с бэкендом
const openFileButton = document.getElementById("openFile");
const fileContentArea = document.getElementById("fileContent");
const passwordInput = document.getElementById("password");
const encryptButton = document.getElementById("encryptFile");
const decryptButton = document.getElementById("decryptFile");
const saveButton = document.getElementById("saveFile");
let filePath;

openFileButton.addEventListener("click", async () => {
  filePath = await window.go.main.App.OpenFile();
  if (filePath) {
    fileContentArea.value = filePath;
  } else {
    alert("Failed to open file.");
  }
});

saveButton.addEventListener("click", async () => {
  const content = fileContentArea.value;
  await window.go.main.App.Save(content);
  alert("File saved!");
});

encryptButton.addEventListener("click", async () => {
  const content = fileContentArea.value;
  const password = passwordInput.value;
  if (!content || !password) {
    alert("Please provide both content and password.");
    return;
  }
  const result = await window.go.main.App.Encrypt(content, password);
  fileContentArea.value = result;
  alert("File encrypted!");
});

decryptButton.addEventListener("click", async () => {
  const content = fileContentArea.value;
  const password = passwordInput.value;
  if (!content || !password) {
    alert("Please provide both content and password.");
    return;
  }
  const result = await window.go.main.App.Decrypt(content, password);
  fileContentArea.value = result;
  alert("File decrypted!");
});
