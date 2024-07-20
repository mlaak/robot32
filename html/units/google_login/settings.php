<?php

$BASE_DIR = __DIR__."/../../..";

// Google OAuth 2.0 credentials
$CLIENT_ID = trim(file_get_contents("$BASE_DIR/keys/google_client_id.txt"));
$CLIENT_SECRET = trim(file_get_contents("$BASE_DIR/keys/google_client_secret.txt"));    
$REDIRECT_URI = trim(file_get_contents("$BASE_DIR/keys/google_redirect_uri.txt")); 
