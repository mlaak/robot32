<?php
require __DIR__."/../middleware.php";

$FAL_KEY = trim(file_get_contents(__DIR__."/../../../keys/falai.txt"));
$url = "https://fal.run/fal-ai/fast-turbo-diffusion";

$content = $_REQUEST["content"];
$data = ["prompt" => $content];
    
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_HTTPHEADER, array(
    "Authorization: Key $FAL_KEY",
    "Content-type: application/json"
));
curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));

$result = curl_exec($ch);
curl_close($ch);

$json = json_decode($result ,true);

//header("Return-Image:".explode(',',$json['images'][0]['url'],2)[1]);

$dat = base64_decode(explode(',',$json['images'][0]['url'],2)[1]);
$md5 = md5($dat);
$time = microtime(true);
$fname = "../../recieved_images/$md5.$time.jpg";

file_put_contents($fname, $dat);

echo "{\"image\":\"$fname\"}";

