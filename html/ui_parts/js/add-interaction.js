function template(name,...replaces){ 
    var s = document.getElementById(name).innerHTML;
    
    
    console.log(replaces);
    
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
      
    cma.insertBefore(message_interaction, cma.firstChild);
    document.getElementById("user-message-"+interaction_no).setAttribute("chat",message);

    
    add_reply_to(cma.firstChild);
    //cma.firstChild.appendChild(reply_to);
    
    return interaction_no;
    
    /*var message_interaction = document.getElementById("template_message_interaction").innerHTML;
    const newHtml = message_interaction.replace("!!USER-REQUEST!!",textToHtml(message)).split("!!REQNO!!").join(interaction_no+"");
    let e = document.getElementById("chat-messages-area");
    //If I have a html code as a string, then I want to make it to html object and insert it after another existing object. How do I do it?
    var newDiv = document.createElement("div");
    newDiv.innerHTML = newHtml;
    
    e.insertBefore(newDiv, e.firstChild);
    
    
    var replyto = document.getElementById("template_reply_to_message").innerHTML;
    const newReplyTo = replyto.split("!!REQNO!!").join(interaction_no+"");
    var newReplyToDiv = document.createElement("div");
    newReplyToDiv.innerHTML = newReplyTo;
    e.firstChild.appendChild(newReplyToDiv);
    
    //e.after(newDiv);
    return interaction_no;*/
}



