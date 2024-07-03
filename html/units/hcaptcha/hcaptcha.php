<?php

function generateRandomString($length = 32) {
    return bin2hex(random_bytes($length / 2));
}


include __DIR__."/settings.php";

$client_id =     trim(file_get_contents("$BASE_DIR/keys/hcaptcha_id.txt"));
$client_secret = trim(file_get_contents("$BASE_DIR/keys/hcaptcha_secret.txt"));    



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
        file_put_contents("$BASE_DIR/working_data/sessions/$session_id.txt","ip, ipus$user_id, -");
        setcookie("r_ression_id",$session_id);
        header('Location: ' . "../index.html");

        echo "verify success";
    } 
    else {
       echo "Verify error - please try again (yes, we know the system is stupid, we are really sorry).";
    }


}
?>



<form method="POST">
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
  <input type="submit" value="Click here">
</form>




<!--
<iframe width = "600" height="400" sandbox="allow-scripts allow-forms allow-same-origin" src="hcaptcha/captcha.html"></iframe>
-->
