<?php

function generateRandomString($length = 32) {
    return bin2hex(random_bytes($length / 2));
}


include __DIR__."/settings.php";

$client_id =     trim(file_get_contents("$BASE_DIR/keys/hcaptcha_id.txt"));
$client_secret = trim(file_get_contents("$BASE_DIR/keys/hcaptcha_secret.txt"));    

$MESSAGE = "";

if(isset($_POST['h-captcha-response'])){

    $data = array(
                'secret' => $client_secret,
                'response' => $_POST['h-captcha-response']
            );
            
    $verify = curl_init();
    curl_setopt($verify, CURLOPT_URL, "https://hcaptcha.com/siteverify");
    curl_setopt($verify, CURLOPT_POST, true);
    curl_setopt($verify, CURLOPT_POSTFIELDS, http_build_query($data));
    curl_setopt($verify, CURLOPT_RETURNTRANSFER, true);
    $response = curl_exec($verify);

     //var_dump($response);
     
     
    $responseData = json_decode($response);

    if($responseData->success) {
        $session_id = generateRandomString(32);
        $headers = getallheaders();

        $user_id = $headers["X-Forwarded-For"];
        header('Create-Session: ' . base64_encode("$session_id, ip, ipus$user_id, -"));
        @file_put_contents("$BASE_DIR/working_data/sessions/$session_id.txt","ip, ipus$user_id, -");
        setcookie("r_ression_id",$session_id,0,"/");
        header('Location: ' . "../index.html");

        $MESSAGE = "verify success"; 
        //echo "verify success";
    } 
    else {
       //echo "Verify error - please try again (yes, we know the system is stupid, we are really sorry).";
       $MESSAGE = "Verify error - please try again (yes, we know the system is stupid, we are really sorry).";
    }


}
?>

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

<p class="text-zinc-800 dark:text-zinc-200">
<?php echo $MESSAGE; ?>
</p>


<center>
<form method="POST" style="margin-top:20px">
  <div class="h-captcha" data-sitekey="<?php echo $client_id;?>"></div>
  <script src="https://js.hcaptcha.com/1/api.js" async defer></script> 
  
  <script>
  function printCookies() {
      const cookies = document.cookie.split(';');
      for (const cookie of cookies) {
        const [name, value] = cookie.split('=');
        console.log(`Name: ${name.trim()}, Value: ${value}`);
      }
    }

    //printCookies();
  </script>
  <input type="submit" value="Click here to continue after proving human" class = "p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full" style="width:400px;margin-top:20px">
</form>
</center>



<!--
<iframe width = "600" height="400" sandbox="allow-scripts allow-forms allow-same-origin" src="hcaptcha/captcha.html"></iframe>
-->
