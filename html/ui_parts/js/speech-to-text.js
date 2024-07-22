// Speech-to-text
const microphoneBtn = document.getElementById('microphone-btn');

window.SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
let recognition = null;
if(window.SpeechRecognition){
    recognition = new SpeechRecognition("sv-SE");
}
else recognition = {};
//recognition.continuous = true; 
recognition.lang = "sv-SE";

function beep() {
  const audio = new Audio(audiofile); // Very short beep sound
  audio.play();
}

microphoneBtn.addEventListener('click', () => {
  if (recognition.continuous) {
    recognition.continuous = false;
    recognition.stop();
    microphoneBtn.classList.remove('bg-red-500', 'text-white');
    microphoneBtn.classList.add('bg-gray-200', 'dark:bg-gray-700');
    //beep();
  } else {
    recognition.continuous = true;
    recognition.start();
    recognition.continuous = true;
    beep();
    microphoneBtn.classList.add('bg-red-500', 'text-white');
    microphoneBtn.classList.remove('bg-gray-200', 'dark:bg-gray-700');
  }
});

recognition.onresult = (event) => {
  const transcript = event.results[event.results.length-1][0].transcript;

  if(transcript.replace(/\s/g, "")=="stopstop"){
    recognition.continuous = false;
    recognition.stop();
    microphoneBtn.classList.remove('bg-red-500', 'text-white');
    microphoneBtn.classList.add('bg-gray-200', 'dark:bg-gray-700');
  }
  else messageBox.value += transcript + ' '; 
};

recognition.onend = () => {
  microphoneBtn.classList.remove('bg-red-500', 'text-white');
  microphoneBtn.classList.add('bg-gray-200', 'dark:bg-gray-700');
  recognition.stop();
    recognition.continuous = false;

  beep();
  run_message(document.getElementById('message-box').value);
  document.getElementById('message-box').value="";
  //recognition = new SpeechRecognition();
};
