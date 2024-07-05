<?php

$GROQ_API_KEY = trim(file_get_contents(__DIR__."/../../../keys/groq.txt"));
$url = "https://api.groq.com/openai/v1/chat/completions";
        
$headers = [
    "Authorization: Bearer $GROQ_API_KEY",
    "Content-Type: application/json"
];


$content = $_GET["content"];
$model = $_GET["model"];

$data = [
    "messages" => [
        [
            "role" => "user",
            "content" => $content
        ]
    ],
    "model" => $model,
    "temperature"=> 1,
    "max_tokens"=> 8024,
    "top_p"=> 1,
    "stream"=> true,
    "stop"=> null

];



$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_POST, 1);
curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
curl_setopt($ch, CURLOPT_RETURNTRANSFER, false);
curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 10);
curl_setopt($ch, CURLOPT_TIMEOUT, 60);

$all_data = "";

$outputted_data = "";

function process_lines($lines){
    global $outputted_data;
    foreach($lines as $l){
        //$json_str = "{".$l."}";
        $json_str = "";
        if(substr($l,0,5)=="data:"){
            $json_str = substr($l,5);
        }
        
        $json = json_decode($json_str,true);
        
        if(isset($json["choices"]) && isset($json["choices"][0]) && isset($json["choices"][0]['delta']) && isset($json["choices"][0]['delta']['content'])){
            $outputted_data.= $json["choices"][0]['delta']['content'];
            echo $json["choices"][0]['delta']['content'];
        }
        
        //echo "\n\n".$json_str."\n\n";
        //print_r($json);
    }

}


curl_setopt($ch, CURLOPT_WRITEFUNCTION, function($curl, $data) {
    global $all_data;
    
    $all_data.=$data;
    
    $lines = explode("\n",$all_data);
    if($lines[count($lines)-1]!=""){
        $all_data = array_pop($lines);
    }
    else {
        array_pop($lines);
        $all_data = "";
    }
   // echo "---------------------------------";
    process_lines($lines);
    flush();
    ob_flush();
    ob_clean();
    // Send the data chunk to the client
  //echo  "\n\n\n\n".$data."\n\n\n\n";
    // Return the length of the data to indicate how much was processed
    return strlen($data);
});


$response = curl_exec($ch);

if (curl_errno($ch)) {
    echo 'Error: ' . curl_error($ch);
} else {
   // echo $response;
}

curl_close($ch);

$currentTime = time();
$year = date('Y', $currentTime);
$month = date('m', $currentTime);
$day = date('d', $currentTime);
$hour = date('H', $currentTime);
$minute = date('i', $currentTime);
$second = date('s', $currentTime);

$time = "$year.$month.$day..$hour.$minute.$second";


$filename = $time . "___" . microtime(true);
$filename = str_replace(".", "_", $filename); // replace the decimal with an underscore
file_put_contents("../../../collected_data/chats/".$filename.".txt", "Model: $model\n\n"."Query:\n".$content."\n\n\nResult:\n".$outputted_data);



