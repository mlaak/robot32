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

<div id="screen" class="flex flex-col md:flex-row h-screen">
    
    <div id="sidebar" class="md:flex-none md:w-48 bg-white dark:bg-zinc-800 shadow-md p-4 space-y-4" >
      <div class="flex justify-between items-center md:block">
        <div class="text-lg font-semibold text-zinc-800 dark:text-zinc-200" id="sidebar-title">Menu</div>
        <button class="md:hidden p-2 bg-zinc-200 dark:bg-zinc-700 text-zinc-800 dark:text-zinc-200 rounded-md focus:outline-none"
          id="sidebar-toggle">
          ☰
        </button>
      </div>     
      <?php echo "<!--";require __DIR__."/sidebar/menu.htm";?>
    </div> <!--sidebar-->


    <div id="main-content" class="flex flex-col flex-1" >  
      <?php echo "<!--";require __DIR__."/topbar/chat-input-area.htm";?> 
      <?php echo "<!--";require __DIR__."/topbar/options.htm";?>  

      <div class="flex-1 overflow-y-auto p-4 space-y-4" id="chat-messages-area">
        <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-md text-lg" id="initial-message">
          <p class="text-zinc-800 dark:text-zinc-200">Hello! How can I assist you today? <button class="p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full md:w-auto" onclick="load_previous_chat();">Load previous</button> </p>
        </div>
        <?php require __DIR__."/maincontent/welcome.htm";?>
      </div>  <!--chat-messages-area-->
      
    </div> <!--main-content-->
    
    
    <div id="parameters-section" class="md:flex-none md:w-48 bg-white dark:bg-zinc-800 shadow-md p-4 space-y-4">
        <?php require __DIR__."/rightbar/parameters.htm";?>  
    </div>
        
</div>  <!--screen-->


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
  
<div id="templates" style="display:none">     
    <?php echo "<!--";require __DIR__."/templates/template_reply_to_message.htm";?>
    <?php echo "<!--";require __DIR__."/templates/template_message_interaction.htm";?>
    <?php echo "<!--";require __DIR__."/templates/template_message_continue.htm";?>
    <?php echo "<!--";require __DIR__."/templates/template_picture_modal.htm";?>

    <?php echo "<!--";require __DIR__."/templates/tpl_ai_response_cont.htm";?>
    <?php echo "<!--";require __DIR__."/templates/tpl_ai_response_initial.htm";?>
    <?php echo "<!--";require __DIR__."/templates/tpl_interaction.htm";?>
    <?php echo "<!--";require __DIR__."/templates/tpl_usr_query_cont.htm";?>
    <?php echo "<!--";require __DIR__."/templates/tpl_usr_query_inital.htm";?>
</div>  
 

<script src="beep.jsi"></script>
<script>

    <?php include __DIR__."/js/functions.js";?>

    if(!checkCookieExists('r_ression_id')){
        window.location.href="login.html";
    }

    var ClientSecretKey=null;
    async function loadSecret(){
        try{
          var cook = getCookie("r_user_secret");
          if(cook!=null){
            var secret = decodeURIComponent(cook);
            ClientSecretKey = await window.crypto.subtle.importKey("jwk",JSON.parse(secret),"AES-GCM",true, ['encrypt', 'decrypt']);
          }
        }
        catch(err){
          console.log(err);
        }
    }
    loadSecret();

    document.getElementById('sidebar-toggle').addEventListener('click', () => {
          document.getElementById('sidebar-menu').classList.toggle('hidden');
    });

    <?php include __DIR__."/js/add-interaction.js"; ?>
    <?php include __DIR__."/js/load-previous-chat.js"; ?> 
    <?php include __DIR__."/js/run-message.js"; ?> 
    <?php include __DIR__."/js/make-modal.js"; ?>  
    <?php include __DIR__."/js/speech-to-text.js"; ?>
</script>

</body>
</html>
