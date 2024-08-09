<?php

function getOpenRouterFreeCredits($apiKey) {
    $ch = curl_init();

    curl_setopt($ch, CURLOPT_URL, 'https://openrouter.ai/api/v1/auth/key');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Authorization: Bearer ' . $apiKey,
        'Content-Type: application/json'
    ]);

    $response = curl_exec($ch);

    if (curl_errno($ch)) {
        throw new Exception('Curl error: ' . curl_error($ch));
    }

    curl_close($ch);

    $data = json_decode($response, true);
    //print_r($data);

    if (isset($data['data']['limit_remaining'])) {
        return $data['data']['limit_remaining'];
    } else {
        throw new Exception('Unable to retrieve free credit balance');
    }
}




function getOpenRouterFreeCreditsTotal($apiKeys) {
    $totalCredits = 0;
    $invalidKeys = [];

    foreach ($apiKeys as $key) {
        try {
            $credits = getOpenRouterFreeCredits($key);
            $totalCredits += $credits;
        } catch (Exception $e) {
            $invalidKeys[] = $key;
        }
    }

    return [
        'totalCredits' => $totalCredits,
        'invalidKeys' => $invalidKeys
    ];
}

/*try {
    $apiKey = $_GET["openrouter_key"];
    $freeCredits = getOpenRouterFreeCredits($apiKey);
    echo $freeCredits;
} catch (Exception $e) {
    echo "Error: " . $e->getMessage();
}*/


if(isset($_REQUEST['r_user_key'])){
    $_COOKIE['r_user_key'] = $_REQUEST['r_user_key'];
}
$OPENROUTER_API_KEY = $_COOKIE['r_user_key'] ?? "";

$arh = apache_request_headers();
$user_keys = $arh["User-Openrouter-Keys"] ?? "";
if($user_keys!=="")$keys = explode(",",$user_keys);
else $keys = [];
if($OPENROUTER_API_KEY!="")array_unshift($keys,$OPENROUTER_API_KEY);
if(count($keys)==0){
    echo "-1";
    exit();
}
else {
    $fc = getOpenRouterFreeCreditsTotal($keys);
    echo $fc["totalCredits"];
    //TODO: handle invalid credits and 0 credits - check and remove them from system
}

