  
  
  
    function run_llm_query(llm_url,message,gpt_text_elem){
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
                //let elem = document.getElementById("gpt-text-"+reqno);
                 const headerString = Array.from(response.headers.entries())
                  .map(([key, value]) => `${key}: ${value}`)
                  .join('\n');
                  
      
                gpt_text_elem.innerHTML = gpt_text_elem.innerHTML+textToHtml(headerString+"");
                
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
                //let elem = document.getElementById("gpt-text-"+reqno);
                gpt_text_elem.innerHTML = gpt_text_elem.innerHTML+textToHtml(decoder.decode(value));
                
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
  
    function run_message_continue(el,el2,message){
        //let reqno = add_interaction(message);
        let reqno = add_interaction_continue(el,el2,message);
        
        let llm_url = document.getElementById("model-select").value;
        gpt_text_elem = document.getElementById("gpt-text-cont-"+reqno);
        
        run_llm_query(llm_url,message,gpt_text_elem);
        
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
        gpt_text_elem = document.getElementById("gpt-text-"+reqno);
        
        run_llm_query(llm_url,message,gpt_text_elem);
        
        
    }
