function get_chat_history(element,out = []){
  for (let i = 0; i < element.children.length; i++) {
    const child = element.children[i];
    if(child.getAttribute('role') == "user" || child.getAttribute('role') == "ai"){
      out.push(child.getAttribute('role')+":"+child.getAttribute('chat'));
    }
    else if(child.tagName=="DIV"){
      get_chat_history(child,out);
    }
  }
  return out;
}

async function saveHistoryImg(parent_req_no,val){
  let encryptedData = await encryptText(val,ClientSecretKey);
  localStorage.setItem("img"+parent_req_no,encryptedData);
}

async function saveHistory(history,parent_req_no){
  let encryptedData = await encryptText(JSON.stringify(history),ClientSecretKey);        
  localStorage.setItem("history"+parent_req_no,encryptedData);
  console.log(localStorage.getItem("history"+parent_req_no));
}
  
function run_llm_query(llm_url,parent_req_no,message,history,gpt_text_elem){
    let history_str = "";
    let llm_response = "";

    //TODO: make this more sane (function should probably take these in as parameters)
    var llm_temp = getSelectBoxValue("llm_temp");
    var llm_top_p = getSelectBoxValue("llm_top_p");
    var llm_max_tokens = getSelectBoxValue("llm_max_tokens");


    if(history!=null){
        history_str = "&history="+encodeURIComponent(JSON.stringify(history));
    }
    fetch(llm_url+'&content='+encodeURIComponent(message)+history_str+"&llm_temp="+llm_temp+"&llm_top_p="+llm_top_p+"&llm_max_tokens="+llm_max_tokens , {credentials: "same-origin"})
      .then(response => {

        checkCredits(document.getElementById("credits_display"));

        if (!response.ok) {
            if(response.status==498){
                window.location.href="login.html";
            }    
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
              return;
            }

            llm_response+=decoder.decode(value);

            gpt_text_elem.innerHTML = gpt_text_elem.innerHTML+textToHtml(decoder.decode(value));
            gpt_text_elem.parentNode.parentNode.setAttribute("chat",gpt_text_elem.parentNode.parentNode.getAttribute("chat")+decoder.decode(value));
            
            return read(); // Read the next chunk
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
    
    if(llm_url!="R32"){
      run_llm_query(llm_url,parent_req_no,message,history,gpt_text_elem);
    }
    else {
      let expert_url = llm_chat.getAttribute("expert_url");
      if(!expert_url)expert_url = "experts/general?model=auto";
      run_llm_query(expert_url,parent_req_no,message,history,gpt_text_elem);
    }
    
}
    
function run_message(message){
    let reqno = add_interaction(message);
    
    localStorage.setItem("reqno",reqno);

    let messageBox = document.getElementById("message-box");
    messageBox.value = "";
    messageBox.style.height = 'auto';
    messageBox.style.height = messageBox.scrollHeight + 'px';
    
    fetch('experts/illustrator?content='+encodeURIComponent(message)+"&illustrator_choice="+getSelectBoxValue("illustrator_choice"),{credentials: "same-origin"})
      .then(response => {
        var imgdata = response.headers.get("Return-Image");
        document.getElementById("gpt-image-"+reqno).src = "data:image/jpeg;base64,"+imgdata;
        return response.json();
      })
      .then(data => {
        document.getElementById("gpt-image-"+reqno).src_link = data.image;
        
        saveHistoryImg(reqno,data.image);

        var picture_modal = document.getElementById("template_picture_modal").innerHTML;
        const newHtml = picture_modal.split("!!REQNO!!").join(interaction_no+"").split("!!IMGSRC!!").join(data.image);
        
        let e = document.getElementById("modal-0");
        const newDiv = document.createElement("div");
        newDiv.innerHTML = newHtml;
        e.parentNode.insertBefore(newDiv, e);
        
        make_modal(reqno);
          
      })
      .catch(error => {
        console.error('Error:', error);
      });
    
    let llm_url = document.getElementById("model-select").value;
    gpt_text_elem = document.getElementById("gpt-text-"+reqno);

    let llm_chat = findParentByClassName(gpt_text_elem,"llm_interaction");
    
    if(llm_url!="R32"){
      run_llm_query(llm_url,reqno,message,null,gpt_text_elem);
    }
    else {
      fetch('experts/classifier?model=mistralai/mixtral-8x7b-instruct&content='+encodeURIComponent(message),{credentials: "same-origin"})
      .then(response => response.text())
      .then(data => {

        if(data.includes("H0")){
          run_llm_query("experts/humorist?a=1",reqno,message,null,gpt_text_elem);
          llm_chat.setAttribute("expert_url","experts/humorist?a=1");
        }
        else if( data.includes("G0") ){

          llm_chat.setAttribute("expert_url","experts/general?model=auto");
          run_llm_query("experts/general?model=auto",reqno,message,null,gpt_text_elem);
        }
        else {
          llm_chat.setAttribute("expert_url","experts/robotics?model=auto");
          run_llm_query("experts/robotics?model=auto",reqno,message,null,gpt_text_elem);
        }

        /*
        if(data.includes("HUM")){
          run_llm_query("experts/humorist?a=1",reqno,message,null,gpt_text_elem);
          llm_chat.setAttribute("expert_url","experts/humorist?a=1");
        }
        else if(data.includes("CCELE") || data.includes("CCPRO") || data.includes("CCAI") || 
          data.includes("CCTEC")    
        ){
          llm_chat.setAttribute("expert_url","experts/robotics?model=auto");
          run_llm_query("experts/robotics?model=auto",reqno,message,null,gpt_text_elem);
        }
        else {
          llm_chat.setAttribute("expert_url","experts/general?model=auto");
          run_llm_query("experts/general?model=auto",reqno,message,null,gpt_text_elem);
        }
        */

      });
    }
}
