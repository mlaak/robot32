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

<div class="justify-center">
    <div class="flex items-center justify-center ">
        <div>
            <!--<div class="w-64 h-64 bg-gray-300 rounded-lg"></div>-->
            
            <div class="flex rounded-md border border-white-500 " style="margin-top:20px"> 
                <img class="rounded-lg border-red-500" src="openscreen.png" width="693" height="535" style="height:535px;widht:693px">  
   
            </div> 
           
            
        </div>            
    </div>   
</div>
 <div class ="flex justify-center" >
                <button class="p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full " style="width:200px;margin:10px" onclick="window.location.href='units/google_login'" > LOGIN WITH GOOGLE </button> 
                
                <button class="p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full " style="width:200px;margin:10px" onclick="window.location.href='units/github_login'"> LOGIN WITH GITHUB</button>
                <button class="p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full " style="width:220px;margin:10px" onclick="window.location.href='units/hcaptcha'"> CONTINUE W/O LOGIN </button>          
            </div>

</body>

<script>
async function generateKey() {
    const key = await window.crypto.subtle.generateKey(
        { name: "AES-GCM", length: 256 },
        true,
        ["encrypt", "decrypt"]
    );
    return key;
}



async function makeKey(){
    const myKey = await generateKey();
    const myKeyStr = JSON.stringify(await window.crypto.subtle.exportKey("jwk", myKey));

    document.cookie = "propesedKey="+myKeyStr;

    checkKey = await window.crypto.subtle.importKey("jwk",JSON.parse(myKeyStr),"AES-GCM",true, ['encrypt', 'decrypt']);

    console.log(myKeyStr);
    console.log(checkKey);
}

makeKey();

</script>
</html>

