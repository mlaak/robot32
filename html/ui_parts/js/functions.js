function textToHtml(text) {
  const entities = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&apos;'
  };
  return text.replace(/[&<>"']/g, (match) => entities[match]);
}


function getCookie(name) {
  const cookieString = document.cookie;
  const cookieStart = cookieString.indexOf(`${name}=`);
  if (cookieStart === -1) {
    return null;
  }
  const cookieEnd = cookieString.indexOf(';', cookieStart);
  if (cookieEnd === -1) {
    return cookieString.substring(cookieStart + name.length + 1);
  }
  return cookieString.substring(cookieStart + name.length + 1, cookieEnd);
}


function checkCookieExists(name) {
  const cookies = document.cookie.split(';');
  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i].trim();
    if (cookie.startsWith(name + '=')) {
      return true;
    }
  }
  return false;
}

function naturalSort(arr) {
  const collator = new Intl.Collator(undefined, {
      numeric: true,
      sensitivity: 'base'
  });
  return arr.sort(collator.compare);
}

async function decryptText(dt,secretKey){
  dt = JSON.parse(dt);
  let decrypted = null;
  try{
    let iva = new Uint8Array(dt.iv);
    let ena = new Uint8Array(dt.encrypted);
  
    decrypted = await window.crypto.subtle.decrypt(
      {name: "AES-GCM",iv: iva},
      secretKey,ena
    );
  }
  catch(e){
    console.log("GOT ERROR from window.crypto.subtle.decrypt",e);
    return null;
  }

  let dec = new TextDecoder();
  return dec.decode(decrypted);
}

async function encryptText(txt,secretKey){
  const textEncoder = new TextEncoder();
  let dt = textEncoder.encode(txt);
  let iv = window.crypto.getRandomValues(new Uint8Array(12)); // the initialization vector should be unique for each encryption
  const encryptedData = await window.crypto.subtle.encrypt(
    {
      name: 'AES-GCM',
      iv: iv
    },
    secretKey,
    dt
  );
  let ee= Array.from(new Uint8Array(encryptedData));
  let decrypted = await window.crypto.subtle.decrypt(
    {name: "AES-GCM",iv: iv},
    secretKey,new Uint8Array(ee)
  );

  return JSON.stringify({"iv":Array.from(iv),"encrypted":Array.from(new Uint8Array(encryptedData))});
}

function getSelectBoxValue(id) {
  let selectBox = document.getElementById(id);
  let selectedValue = selectBox.options[selectBox.selectedIndex].value;
  return selectedValue;
}

function deleteAllLLMInteractions(){
  var elements = document.querySelectorAll('.llm_interaction');
  // Loop through the elements and remove them from the DOM
  for (var i = 0; i < elements.length; i++) {
    elements[i].remove();
  }
}