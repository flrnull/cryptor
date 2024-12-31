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

const openFileButton = document.getElementById("openFile");
const fileContentArea = document.getElementById("fileContent");
const passwordInput = document.getElementById("password");
const encryptButton = document.getElementById("encryptFile");
const decryptButton = document.getElementById("decryptFile");
let filePath;

// Функция для показа диалогов ошибок
async function showError(message) {
  await window.runtime.MessageDialog({
    type: "error",
    title: "Error",
    message: message,
  });
}

openFileButton.addEventListener("click", async () => {
  try {
    filePath = await window.go.main.App.OpenFile();
    if (filePath) {
      fileContentArea.value = filePath;
    } else {
      await showError("Failed to open file.");
    }
  } catch (error) {
    await showError(`Error: ${error.message}`);
  }
});

encryptButton.addEventListener("click", async () => {
  const content = fileContentArea.value;
  const password = passwordInput.value;
  if (!content || !password) {
    await showError("Please provide both content and password.");
    return;
  }
  try {
    const result = await window.go.main.App.Encrypt(content, password);
    fileContentArea.value = result;
    await window.runtime.MessageDialog({
      type: "info",
      title: "Success",
      message: "File encrypted successfully!",
    });
  } catch (error) {
    await showError(`Encryption Error: ${error.message}`);
  }
});

decryptButton.addEventListener("click", async () => {
  const content = fileContentArea.value;
  const password = passwordInput.value;
  if (!content || !password) {
    await showError("Please provide both content and password.");
    return;
  }
  try {
    const result = await window.go.main.App.Decrypt(content, password);
    fileContentArea.value = result;
    await window.runtime.MessageDialog({
      type: "info",
      title: "Success",
      message: "File decrypted successfully!",
    });
  } catch (error) {
    await showError(`Decryption Error: ${error.message}`);
  }
});