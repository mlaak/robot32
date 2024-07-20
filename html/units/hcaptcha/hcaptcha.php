<?php
include __DIR__."/settings.php";

function generateRandomString($length = 32) {
    return bin2hex(random_bytes($length / 2));
}
    
$MESSAGE = "";
if(isset($_POST['h-captcha-response'])){
    $data = [
            'secret' => $CLIENT_SECRET,
            'response' => $_POST['h-captcha-response']
            ];
    $verify = curl_init();
    curl_setopt($verify, CURLOPT_URL, "https://hcaptcha.com/siteverify");
    curl_setopt($verify, CURLOPT_POST, true);
    curl_setopt($verify, CURLOPT_POSTFIELDS, http_build_query($data));
    curl_setopt($verify, CURLOPT_RETURNTRANSFER, true);
    $response = curl_exec($verify);     
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
    } 
    else {
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
  <div class="h-captcha" data-sitekey="<?php echo $CLIENT_ID;?>"></div>
  <script src="https://js.hcaptcha.com/1/api.js" async defer></script> 
  <input type="submit" value="Click here to continue after proving human" class = "p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full" style="width:400px;margin-top:20px">
</form>
</center>
</body>
</html>
