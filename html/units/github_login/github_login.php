<?php
//TODO: This file is a bit of a mess
include __DIR__."/settings.php";

function generateRandomString($length = 32) {
    return bin2hex(random_bytes($length / 2));
}


// GitHub OAuth settings
$client_id =     trim(file_get_contents("$BASE_DIR/keys/github_client_id.txt"));
$client_secret = trim(file_get_contents("$BASE_DIR/keys/github_client_secret.txt"));    
$redirect_uri =  trim(file_get_contents("$BASE_DIR/keys/github_redirect_uri.txt"));    


// Authorization
if (!isset($_GET['code'])) {
    $auth_url = "https://github.com/login/oauth/authorize?client_id=$client_id&redirect_uri=$redirect_uri&scope=user:email";
    header("Location: $auth_url");
    exit;
}

// Exchange code for access token
$code = $_GET['code'];
$token_url = 'https://github.com/login/oauth/access_token';
$data = [
    'client_id' => $client_id,
    'client_secret' => $client_secret,
    'code' => $code,
    'redirect_uri' => $redirect_uri
];

$curl = curl_init($token_url);
curl_setopt($curl, CURLOPT_POST, true);
curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query($data));
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
curl_setopt($curl, CURLOPT_HTTPHEADER, ['Accept: application/json']);

$response = curl_exec($curl);
curl_close($curl);

$token_data = json_decode($response, true);
$access_token = $token_data['access_token'];

// Get user data
$api_url = 'https://api.github.com/user';
$curl = curl_init($api_url);
curl_setopt($curl, CURLOPT_HTTPHEADER, [
    'Authorization: token ' . $access_token,
    'User-Agent: PHP Script'
]);
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);

$response = curl_exec($curl);
curl_close($curl);

$user_data = json_decode($response, true);
$user_id = $user_data['id'];


// Get user email
$email_url = 'https://api.github.com/user/emails';
$curl = curl_init($email_url);
curl_setopt($curl, CURLOPT_HTTPHEADER, [
    'Authorization: token ' . $access_token,
    'User-Agent: PHP Script'
]);
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);

$response = curl_exec($curl);
curl_close($curl);

$email_data = json_decode($response, true);
$primary_email = '';

foreach ($email_data as $email) {
    if ($email['primary']) {
        $primary_email = $email['email'];
        break;
    }
}

if($primary_email!=""){
    $session_id = generateRandomString(32);
    $user_id = ($user_data['id']*1)."";
    if(!is_numeric($user_id))exit("USER ID SHOULD BE A NUMBER");
    $user_email = $primary_email;
    
    header('Create-Session: ' . base64_encode("$session_id, github, gith$user_id, $user_email"));
    @file_put_contents("$BASE_DIR/working_data/sessions/$session_id.txt","github, gith$user_id, $user_email");
    
    $user_id = "gith$user_id";
    if(!file_exists("$BASE_DIR/working_data/users/$user_id/user_id.txt" )){
        $secret_key = $_COOKIE["propesedKey"];
        
        $userdir = "$BASE_DIR/working_data/users/$user_id";
        @mkdir($userdir,0777,true);
        file_put_contents("$userdir/user_id.txt",$user_id);
        file_put_contents("$userdir/user_email.txt",$user_email);
        file_put_contents("$userdir/secret_key.txt",$secret_key);     
        setcookie("r_user_secret",$secret_key,0,"/"); 
    } else {
        if(file_exists("$BASE_DIR/working_data/users/$user_id/secret_key.txt" )){
            $userdir = "$BASE_DIR/working_data/users/$user_id";
            $secret_key = file_get_contents("$userdir/secret_key.txt");
            setcookie("r_user_secret",$secret_key,0,"/");
        }
    }
    
    setcookie("propesedKey","",0,"/");
    setcookie("r_ression_id",$session_id,0,"/");
    header('Location: ' . "../index.html");
}



