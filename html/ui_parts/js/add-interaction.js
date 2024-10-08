function template(name,...replaces){ 
    var s = document.getElementById(name).innerHTML;
    
    if(replaces){
        replaces.forEach(function(r) {
            s = s.split(r[0]).join(r[1]);
        });
    }
    var newDiv = document.createElement("div");
    newDiv.innerHTML = s;
    return newDiv;
}

function findParentByClassName(element, className) {
    let currentElement = element;
  
    if (currentElement.classList.contains(className)) {
        return currentElement;
    }

    while (currentElement && currentElement.parentNode) {
      currentElement = currentElement.parentNode;
  
      if (currentElement.classList.contains(className)) {
        return currentElement;
      }
    }
    return null; // return null if not found
}



let interaction_no = 0;

interaction_no = localStorage.getItem("reqno")*1;

function add_interaction_continue(el,el2,message){
    interaction_no += 1;

    el.removeChild(el2);    
     let message_interaction = template("template_message_continue",
                                ["!!USER-REQUEST!!",textToHtml(message)],
                                ["!!REQNO!!",interaction_no+""]);   
     el.appendChild(message_interaction);      
     document.getElementById("user-message-cont-"+interaction_no).setAttribute("chat",message);
     add_reply_to(el);
     return interaction_no;               
    
}


function add_reply_to(el){
    let reply_to = template("template_reply_to_message", ["!!REQNO!!",interaction_no+""]);
    el.appendChild(reply_to);
}

function add_interaction(message){
    interaction_no += 1;    
    let cma = document.getElementById("chat-messages-area");

    let message_interaction = template("template_message_interaction",
                                ["!!USER-REQUEST!!",textToHtml(message)],
                                ["!!REQNO!!",interaction_no+""]); 

    message_interaction.classList.add("llm_interaction");
    message_interaction.setAttribute("interaction_no",interaction_no);
      
    cma.insertBefore(message_interaction, cma.firstChild);
    document.getElementById("user-message-"+interaction_no).setAttribute("chat",message);
    
    add_reply_to(cma.firstChild);
    
    return interaction_no;
    
}



function restore_interaction(reqno,img,history){
    if(history == null)return 0;
    if(document.getElementById("interaction-"+reqno)!=null)return 0;

    let cma = document.getElementById("chat-messages-area");
        
    let message_interaction = template("tpl_interaction",
                                ["!!REQNO!!",reqno]); 

    message_interaction.classList.add("llm_interaction");
    message_interaction.setAttribute("interaction_no",reqno);
      
    cma.insertBefore(message_interaction, document.getElementById("initial-message"));
    
    let firstUser = true;
    let firtResponse = true;

    for(i=0;i<history.length;i++){
        let h = history[i]+"";
        if(h.startsWith("user:")){
            let mess = h.substring(5);
            interaction_no += 1;

            if(firstUser){
                var usrt = template("tpl_usr_query_initial", 
                    ["!!USER-REQUEST!!",textToHtml(mess)],
                    ["!!REQNO!!",interaction_no+""]);
                firstUser = false;
            } else {
                var usrt = template("tpl_usr_query_cont", 
                    ["!!USER-REQUEST!!",textToHtml(mess)],
                    ["!!REQNO!!",interaction_no+""]);
            }

            message_interaction.appendChild(usrt);
        }
        else if(h.startsWith("ai:")){
            let mess = h.substring(3);
            interaction_no += 1;
            if(firtResponse){
                var airesp = template("tpl_ai_response_initial",
                    ["!!IMG-LINK!!",img], 
                    ["!!AI-RESPONSE!!",textToHtml(mess)],
                    ["!!REQNO!!",interaction_no+""]);
                firtResponse = false;
            }
            else {
                var airesp = template("tpl_ai_response_cont", 
                    ["!!AI-RESPONSE!!",textToHtml(mess)],
                    ["!!REQNO!!",interaction_no+""]);
            }
            message_interaction.appendChild(airesp);
        }
    }
    
    return interaction_no;
}




