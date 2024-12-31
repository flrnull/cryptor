import logo from './assets/images/logo-universal.png';

document.querySelector('#app').innerHTML = `
  <div>
    <img src="${logo}" alt="Logo" />
    <h1>Welcome to Cryptor!</h1>
    <p>Use the app to encrypt or decrypt your files.</p>
  </div>
`;
