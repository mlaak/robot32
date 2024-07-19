
function get_chat_history(element,out = []){
  for (let i = 0; i < element.children.length; i++) {
    const child = element.children[i];
    //console.log("child is",child,child.tagName);
    if(child.getAttribute('role') == "user" || child.getAttribute('role') == "ai"){
      out.push(child.getAttribute('role')+":"+child.getAttribute('chat'));
    }
    else if(child.tagName=="DIV"){
      get_chat_history(child,out);
    }
    // Do something with the child element
    
  }
  return out;
}
    async function saveHistoryImg(parent_req_no,val){
      let encryptedData = await encryptText(val,ClientSecretKey);
      localStorage.setItem("img"+parent_req_no,encryptedData);
    }

    async function saveHistory(history,parent_req_no){
     /* const textEncoder = new TextEncoder();
      let dt = textEncoder.encode(JSON.stringify(history));
      const encryptedData = await window.crypto.subtle.encrypt(
        {
          name: 'AES-GCM',
          iv: window.crypto.getRandomValues(new Uint8Array(12)), // the initialization vector should be unique for each encryption
        },
        ClientSecretKey,
        dt
      );*/
      let encryptedData = await encryptText(JSON.stringify(history),ClientSecretKey);
     
      console.log(encryptedData);
        
      localStorage.setItem("history"+parent_req_no,encryptedData);

      console.log(localStorage.getItem("history"+parent_req_no));



    }
  
    function run_llm_query(llm_url,parent_req_no,message,history,gpt_text_elem){
        let history_str = "";

        let llm_response = "";

        if(history!=null){
            history_str = "&history="+encodeURIComponent(JSON.stringify(history));
        }
        fetch(llm_url+'&content='+encodeURIComponent(message)+history_str,{credentials: "same-origin"})
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
                  if(history==null)history = [];
                  history.push("user:"+message);
                  history.push("ai:"+llm_response);
                  

                  saveHistory(history,parent_req_no);


                  console.log(history);

                  console.log('Stream reading complete');

                  return;
                }

                // Decode the chunk to a string and log it
                //let elem = document.getElementById("gpt-text-"+reqno);
                llm_response+=decoder.decode(value);

                gpt_text_elem.innerHTML = gpt_text_elem.innerHTML+textToHtml(decoder.decode(value));
                gpt_text_elem.parentNode.parentNode.setAttribute("chat",gpt_text_elem.parentNode.parentNode.getAttribute("chat")+decoder.decode(value));
                //console.log(decoder.decode(value));

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
        
        llm_chat = findParentByClassName(el,"llm_interaction");
        let history = get_chat_history(llm_chat);
        
        let reqno = add_interaction_continue(el,el2,message);
        

        var prn = findParentByClassName(el,"llm_interaction");
        let parent_req_no = prn.getAttribute("interaction_no");

        localStorage.setItem("reqno",reqno);

        let llm_url = document.getElementById("model-select").value;
        gpt_text_elem = document.getElementById("gpt-text-cont-"+reqno);
        
        

        run_llm_query(llm_url,parent_req_no,message,history,gpt_text_elem);
        
    }
    
    function run_message(message){
        let reqno = add_interaction(message);
        
        localStorage.setItem("reqno",reqno);

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
            
            saveHistoryImg(reqno,data.image);
            //localStorage.setItem("img"+reqno,data.image);

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
        
        
        if(llm_url!="R32"){
          run_llm_query(llm_url,reqno,message,null,gpt_text_elem);
        }
        else {
          fetch('experts/classifier?model=mistralai/mixtral-8x7b-instruct&content='+encodeURIComponent(message),{credentials: "same-origin"})
          .then(response => response.text())
          .then(data => {
            if(data.includes("CCL1")){
              run_llm_query("experts/humorist?a=1",reqno,message,null,gpt_text_elem);
            }
            else if(data.includes("CCL3.2") || data.includes("CCL4") || data.includes("CCL5") || 
              data.includes("CCL6") || data.includes("CCL7")    
            ){
              run_llm_query("experts/robotics?model=auto",reqno,message,null,gpt_text_elem);
            }
            else {
              run_llm_query("experts/general?model=auto",reqno,message,null,gpt_text_elem);
            }
            // do something with the data here
          });
  //            document.getElementById("gpt-image-"+reqno).src = data.image;
            //document.getElementById("gpt-image-"+reqno).src_link = data.image;
          
        }

        


        
        
        
    }
