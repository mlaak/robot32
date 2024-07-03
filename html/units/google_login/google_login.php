<?php

include __DIR__."/settings.php";

function generateRandomString($length = 32) {
    return bin2hex(random_bytes($length / 2));
}


// Google OAuth 2.0 credentials
$client_id = trim(file_get_contents("$BASE_DIR/keys/google_client_id.txt"));
$client_secret = trim(file_get_contents("$BASE_DIR/keys/google_client_secret.txt"));    
$redirect_uri = trim(file_get_contents("$BASE_DIR/keys/google_redirect_uri.txt")); 
//$redirect_uri = 'http://localhost:8000/google_login.php';

// Google OAuth 2.0 endpoints
$auth_url = 'https://accounts.google.com/o/oauth2/auth';
$token_url = 'https://accounts.google.com/o/oauth2/token';
$userinfo_url = 'https://www.googleapis.com/oauth2/v3/userinfo';

// Step 1: Redirect to Google's authorization page
if (!isset($_GET['code'])) {
    $params = array(
        'client_id' => $client_id,
        'redirect_uri' => $redirect_uri,
        'response_type' => 'code',
        'scope' => 'https://www.googleapis.com/auth/userinfo.email',
        'access_type' => 'online'
    );
    
    $auth_link = $auth_url . '?' . http_build_query($params);
    header('Location: ' . $auth_link);
    
    
    //echo ('Location: ' . $auth_link);
    
    exit;
}

// Step 2: Exchange authorization code for access token
if (isset($_GET['code'])) {
    $code = $_GET['code'];
    
    $curl = curl_init();
    curl_setopt($curl, CURLOPT_URL, $token_url);
    curl_setopt($curl, CURLOPT_POST, true);
    curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query([
        'code' => $code,
        'client_id' => $client_id,
        'client_secret' => $client_secret,
        'redirect_uri' => $redirect_uri,
        'grant_type' => 'authorization_code'
    ]));
    curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
    $response = curl_exec($curl);
    curl_close($curl);
    
    $token_data = json_decode($response, true);
    
    if (isset($token_data['access_token'])) {
        $access_token = $token_data['access_token'];
        
        // Step 3: Fetch user information
        $curl = curl_init();
        curl_setopt($curl, CURLOPT_URL, $userinfo_url);
        curl_setopt($curl, CURLOPT_HTTPHEADER, ['Authorization: Bearer ' . $access_token]);
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
        $response = curl_exec($curl);
        curl_close($curl);
        
        $user_data = json_decode($response, true);
        print_r($user_data);
        if (isset($user_data['email'])) {
            
            $session_id = generateRandomString(32);
            $user_id = $user_data['sub'];
            $user_email = $user_data['email'];
            
            file_put_contents("$BASE_DIR/working_data/sessions/$session_id.txt","google, goog$user_id, $user_email");
            setcookie("r_ression_id",$session_id);
            header('Location: ' . "index.html");
        } else {
            echo "Failed to fetch user information.";
        }
    } else {
        echo "Failed to obtain access token.";
    }
}
?>
