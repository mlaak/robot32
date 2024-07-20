<?php

$BASE_DIR = __DIR__."/../../..";

// GitHub OAuth settings
$CLIENT_ID =     trim(file_get_contents("$BASE_DIR/keys/github_client_id.txt"));
$CLIENT_SECRET = trim(file_get_contents("$BASE_DIR/keys/github_client_secret.txt"));    
$REDIRECT_URI =  trim(file_get_contents("$BASE_DIR/keys/github_redirect_uri.txt"));    

