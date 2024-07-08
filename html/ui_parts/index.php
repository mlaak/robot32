<!DOCTYPE html>
<html lang="en">
<head>
  <script>
  </script>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Chat Interface</title>
  <script src="https://cdn.tailwindcss.com/3.4.3"></script>
</head>
<body class="bg-zinc-100 dark:bg-zinc-900">
<?php
//echo $_SERVER['REMOTE_ADDR'];
?>

  <div class="flex flex-col md:flex-row h-screen">
    
    <div class="md:flex-none md:w-48 bg-white dark:bg-zinc-800 shadow-md p-4 space-y-4" id="sidebar">
      <div class="flex justify-between items-center md:block">
        <div class="text-lg font-semibold text-zinc-800 dark:text-zinc-200" id="sidebar-title">Menu</div>
        <button class="md:hidden p-2 bg-zinc-200 dark:bg-zinc-700 text-zinc-800 dark:text-zinc-200 rounded-md focus:outline-none"
          id="sidebar-toggle">
          ☰
        </button>
      </div>
      
      <?php echo "<!--";require __DIR__."/sidebar/menu.htm";?>
    </div>

    <div class="flex flex-col flex-1" id="main-content">
      
      <?php echo "<!--";require __DIR__."/topbar/chat-input-area.htm";?> 
      <?php echo "<!--";require __DIR__."/topbar/options.htm";?>  


      <div class="flex-1 overflow-y-auto p-4 space-y-4" id="chat-messages-area">


        <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-md text-lg" id="initial-message">
          <p class="text-zinc-800 dark:text-zinc-200">Hello! How can I assist you today?</p>
        </div>


        <div> 
          <div class="bg-blue-500 text-white p-4 rounded-lg shadow-md self-end text-lg " id="user-message">
            <p>What is robot32-chat and what can I do with it.</p>
          </div>

          <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-md flex flex-col sm:flex-row items-start space-y-4 sm:space-y-0 sm:space-x-4 text-lg" id="gpt-message">  
            <img src="/welcome2.jpg" width="200" height="200" alt="generated image"
              class="rounded-md shadow-md cursor-pointer" id="gpt-image-0" />
              <pre style = "white-space: pre-wrap" ><p class="text-zinc-800 dark:text-zinc-200">Robot32-chat is large language model based on open source LLM. 
Type in your question (or picture idea) and hit the RUN button.
You can try:
  What is the capital of Sweden?
  What is the smallest whale?
  Make me pythonc function that sorts custom objects be given param.
  A serene forest landscape at sunset, glassy lake reflecting the sky.
  A surreal, dreamlike landscape with floating islands.
          
Image generation is based on open source model Stable Diffusion Turbo (v1.5/XL)
            </pre>
            </p>
          </div>
        
          
          
          
        </div>
        
        
        
      </div>
    </div>

    <div class="md:flex-none md:w-48 bg-white dark:bg-zinc-800 shadow-md p-4 space-y-4" id="parameters-section">
      <div class="flex justify-center items-center cursor-pointer" onclick="toggleParameters()">
        <div style="color:gray" class="text-lg font-semibold text-zinc-800 dark:text-zinc-200">Parameters</div>
        <svg class="w-4 h-4 ml-1 transform transition duration-200" id="parameters-arrow" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
        </svg>
      </div>
      <ul class="hidden space-y-2" id="parameters-list">
        <li>Parameter 1</li>
        <li>Parameter 2</li>
        <li>Parameter 3</li>
      </ul>
    </div>
  </div>

  <div id="modal-0" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center hidden">
    <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-md space-y-4" id="modal-content-0">
      <img src="welcome2.jpg" alt="generated image" class="w-full h-auto" id="modal-image-0" />
      <div class="flex space-x-4" id="modal-buttons-0">
        <button
          class="p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
          id="download-btn">Download</button>
        <button
          class="p-2 bg-green-500 text-white rounded-md shadow-md hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          id="animate-btn">Animate</button>
      </div>
    </div>
  </div>
  
  
<?php include __DIR__."/templates/template_message_interaction.htm";?>
  
<?php include __DIR__."/templates/template_picture_modal.htm";?>
  
 
  <script>
    function toggleParameters() {
      const parametersList = document.getElementById('parameters-list');
      const arrow = document.getElementById('parameters-arrow');

      parametersList.classList.toggle('hidden');
      arrow.classList.toggle('rotate-180'); 
    }
  </script>



  <script src="beep.jsi"></script>
  <script>
  
  
    /*let message_interaction = `
    <div>
      <div class="bg-blue-500 text-white p-4 rounded-lg shadow-md self-end text-lg" id="user-message-!!REQNO!!">
        <p>!!USER-REQUEST!!</p>
      </div>
      <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-md flex flex-col sm:flex-row items-start space-y-4 sm:space-y-0 sm:space-x-4 text-lg" id="gpt-message--!!REQNO!!"> 
      
    
      
<!--      <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-md flex items-start space-x-4 text-lg"
        id="gpt-message-!!REQNO!!">-->
        <img src="https://placehold.co/200x200" width="200" height="200" alt="generated image"
          class="rounded-md shadow-md cursor-pointer" id="gpt-image-!!REQNO!!" />
        
        <pre style = "white-space: pre-wrap" text-zinc-800 dark:text-zinc-200 ><p id="gpt-text-!!REQNO!!" class="text-zinc-800 dark:text-zinc-200"></p>
        </pre>  
          
       
      </div>   
      
      
      <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-md flex flex-col sm:flex-row items-start space-y-4 sm:space-y-0 sm:space-x-4 text-lg" id="reply-message-div-!!REQNO!!">
      
      
      <textarea style="min-height: 45px; height: 45px;" id="reply-message-box-!!REQNO!!" placeholder="Type a reply..." class="w-full p-2 border border-zinc-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-zinc-800 dark:border-zinc-600 dark:placeholder-zinc-600 dark:text-white pr-12 text-lg resize-none overflow-hidden" run_button="run-button--!!REQNO!!"></textarea> 
      
      <button style="height: 45px; width: 80px"  class="p-2 bg-white dark:bg-zinc-800 text-zinc-400  dark:text-zinc-600 focus:text-white  border border-zinc-300 dark:border-zinc-600 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full" id="run-button--!!REQNO!!" onclick="run_message(document.getElementById('reply-message-box-!!REQNO!!').value);document.getElementById('reply-message-box-!!REQNO!!').value='';">RUN</button>
      
      
      </div>
      
    </div>
  `;
  
  */
  
  
  /*
  let picture_modal = `
  <div id="modal-!!REQNO!!" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center hidden">
    <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-md space-y-4" id="modal-content-!!REQNO!!">
      <img src="!!IMGSRC!!" alt="generated image" class="w-full h-auto" id="modal-image-!!REQNO!!" />
      <div class="flex space-x-4" id="modal-buttons-!!REQNO!!">
        <button
          class="p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
          id="download-btn">Download</button>
        <button
          class="p-2 bg-green-500 text-white rounded-md shadow-md hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          id="animate-btn">Animate</button>
      </div>
    </div>
  </div>
  `;
  */
  
  
  
    let interaction_no = 0;
    
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

    
    function add_interaction(message){
        interaction_no += 1;
        
        var message_interaction = document.getElementById("template_message_interaction").innerHTML;
        const newHtml = message_interaction.replace("!!USER-REQUEST!!",textToHtml(message)).split("!!REQNO!!").join(interaction_no+"");
        let e = document.getElementById("chat-messages-area");
        //If I have a html code as a string, then I want to make it to html object and insert it after another existing object. How do I do it?
        const newDiv = document.createElement("div");
        newDiv.innerHTML = newHtml;
        e.insertBefore(newDiv, e.firstChild);
        //e.after(newDiv);
        return interaction_no;
    }
  
  
  
  
    function run_message(message){
        let reqno = add_interaction(message);
        
        let messageBox = document.getElementById("message-box");
        messageBox.value = "";
        messageBox.style.height = 'auto';
        messageBox.style.height = messageBox.scrollHeight + 'px';
        
        
        
        //fetch('routers/falai/falai.php?content='+encodeURIComponent(message))
        fetch('experts/illustrator?content='+encodeURIComponent(message),{credentials: "same-origin"})
          .then(response => {
            var imgdata = response.headers.get("Return-Image");
            document.getElementById("gpt-image-"+reqno).src = "data:image/jpeg;base64,"+imgdata;
            return response.json();
          })
          .then(data => {
//            document.getElementById("gpt-image-"+reqno).src = data.image;
            document.getElementById("gpt-image-"+reqno).src_link = data.image;
            
            console.log(data);
            
            var picture_modal = document.getElementById("template_picture_modal").innerHTML;
            const newHtml = picture_modal.split("!!REQNO!!").join(interaction_no+"").split("!!IMGSRC!!").join(data.image);
            
            let e = document.getElementById("modal-0");
            const newDiv = document.createElement("div");
            newDiv.innerHTML = newHtml;
            e.parentNode.insertBefore(newDiv, e);
            
            
            make_modal(reqno);
            //modal
            
            
            
            
          })
          .catch(error => {
            console.error('Error:', error);
          });
        
        let llm_url = document.getElementById("model-select").value;
        //fetch('groq.php?content='+message)
        
        
        fetch(llm_url+'&content='+encodeURIComponent(message),{credentials: "same-origin"})
        //fetch('tes.php?content='+message)
          .then(response => {
          //console.log(response)
            if (!response.ok) {
                if(response.status==498){
                    window.location.href="login.html";
                }    
                console.log(response)
                 //const decoder = new TextDecoder();
                let elem = document.getElementById("gpt-text-"+reqno);
                 const headerString = Array.from(response.headers.entries())
                  .map(([key, value]) => `${key}: ${value}`)
                  .join('\n');
                  
      
                elem.innerHTML = elem.innerHTML+textToHtml(headerString+"");
                
              throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.body;
          })
          .then(stream => {
            const reader = stream.getReader();
            const decoder = new TextDecoder();

            // Read the stream in chunks
            function read() {
              return reader.read().then(({ done, value }) => {
                if (done) {
                  console.log('Stream reading complete');
                  return;
                }

                // Decode the chunk to a string and log it
                let elem = document.getElementById("gpt-text-"+reqno);
                elem.innerHTML = elem.innerHTML+textToHtml(decoder.decode(value));
                
                console.log(decoder.decode(value));

                // Read the next chunk
                return read();
              });
            }

            // Start reading the stream
            return read();
          })
          .catch(err => console.error(err));
    }
    
      document.getElementById('sidebar-toggle').addEventListener('click', () => {
          document.getElementById('sidebar-menu').classList.toggle('hidden');
        });
    
    
    function make_modal(reqno){
   
      

        document.getElementById('gpt-image-'+reqno).addEventListener('click', () => {
            console.log("clicked");
          document.getElementById('modal-'+reqno).classList.remove('hidden');
        });

        document.getElementById('download-btn').addEventListener('click', () => {
          const link = document.createElement('a');
          link.href = document.getElementById('modal-image-'+reqno).src;
          link.download = 'generated image';
          link.click();
        });

        document.getElementById('animate-btn').addEventListener('click', () => {
          const image = document.getElementById('modal-image-'+reqno);
          image.classList.add('animate-bounce');
          setTimeout(() => image.classList.remove('animate-bounce'), 1000);
        });

        document.getElementById('modal-'+reqno).addEventListener('click', (e) => {
          if (e.target === document.getElementById('modal-'+reqno)) {
            document.getElementById('modal-'+reqno).classList.add('hidden');
          }
        });
    
    
    }
    make_modal(0);

    
    function checkCookie(name) {
      const cookies = document.cookie.split(';');
      for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.startsWith(name + '=')) {
          return true;
        }
      }
      return false;
    }
    
    if(!checkCookie('r_ression_id')){
        window.location.href="login.html";
    }


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
      console.log(event.results);
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
  </script>
</body>
</html>
