<?php
//TODO: This file is a bit of a mess
include __DIR__."/settings.php";

function generateRandomString($length = 32) {
    return bin2hex(random_bytes($length / 2));
}


// *** Sending the user to github site so they can return with an auth code ***

if (!isset($_GET['code'])) {
    $auth_url = "https://github.com/login/oauth/authorize?client_id=$CLIENT_ID&redirect_uri=$REDIRECT_URI&scope=user:email";
    header("Location: $auth_url");
    exit;
}

// ******** Exchange auth code for access token ********************

$curl = curl_init('https://github.com/login/oauth/access_token');
curl_setopt($curl, CURLOPT_POST, true);
curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query([
            'client_id' => $CLIENT_ID,
            'client_secret' => $CLIENT_SECRET,
            'code' => $_GET['code'],
            'redirect_uri' => $REDIRECT_URI
        ]));
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
curl_setopt($curl, CURLOPT_HTTPHEADER, ['Accept: application/json']);
$response = curl_exec($curl);   curl_close($curl);

$token_data = json_decode($response, true);
$access_token = $token_data['access_token'];


// ***************** Get user data *****************************

$curl = curl_init('https://api.github.com/user');
curl_setopt($curl, CURLOPT_HTTPHEADER, [
    'Authorization: token ' . $access_token,
    'User-Agent: PHP Script'
]);
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
$response = curl_exec($curl);   curl_close($curl);

$user_data = json_decode($response, true);
$user_id = $user_data['id'];


// ***************** Get user email *****************************

$curl = curl_init('https://api.github.com/user/emails');
curl_setopt($curl, CURLOPT_HTTPHEADER, [
    'Authorization: token ' . $access_token,
    'User-Agent: PHP Script'
]);
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
$response = curl_exec($curl);   curl_close($curl);

$email_data = json_decode($response, true);
$primary_email = '';

foreach ($email_data as $e) if($e['primary']) {
    $primary_email = $e['email'];   break;
}

// ***************** IF success then create login *****************************

if($primary_email!=""){
    //user data
    $session_id = generateRandomString(32);
    $user_id = ($user_data['id']*1)."";
    if(!is_numeric($user_id))exit("USER ID SHOULD BE A NUMBER");
    $user_email = $primary_email;
    
    //create session for the user
    header('Create-Session: ' . base64_encode("$session_id, github, gith$user_id, $user_email"));
    @file_put_contents("$BASE_DIR/working_data/sessions/$session_id.txt","github, gith$user_id, $user_email");
    

    $user_id = "gith$user_id";
    if(!file_exists("$BASE_DIR/working_data/users/$user_id/user_id.txt" )){
        //create user 
        $secret_key = $_COOKIE["propesedKey"];
        
        $userdir = "$BASE_DIR/working_data/users/$user_id";
        @mkdir($userdir,0777,true);
        file_put_contents("$userdir/user_id.txt",$user_id);
        file_put_contents("$userdir/user_email.txt",$user_email);
        file_put_contents("$userdir/secret_key.txt",$secret_key);     
        setcookie("r_user_secret",$secret_key,0,"/"); 
    } else {
        //load the secret from existing user
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



