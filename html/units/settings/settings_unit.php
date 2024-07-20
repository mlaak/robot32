<?php
$MESSAGE = "";
if(isset($_REQUEST['code'])){
    $code = $_GET['code']; 
    $data = ['code' => $code];

    $options = [
      'http' => [
        'header'  => "Content-type: application/json\r\n",
        'method'  => 'POST',
        'content' => json_encode($data),
      ],
    ];
    $context  = stream_context_create($options);
    $result = file_get_contents('https://openrouter.ai/api/v1/auth/keys', false, $context);

    if ($result === FALSE) { $MESSAGE = "Unsuccessful"; }
    else {
        $data = json_decode($result,TRUE);
        setcookie('r_user_key', $data['key'], 0, '/'); 
        $MESSAGE = "Successfully connected. ";
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

  <script>
    theurl = "https://openrouter.ai/auth?callback_url="+window.location.href;
  </script>

  <center>
    <p class="text-zinc-800 dark:text-zinc-200">
      <button 
        class = "p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full" 
        style="width:400px;margin-top:20px" 
        onclick="window.location.href=theurl;">
        CONNECT OPENROUTER
      </button>
      <br>
      <button  
        class = "p-2 bg-blue-500 text-white rounded-md shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg w-full" 
        style="width:400px;margin-top:20px" 
        onclick="window.location.href='../index.html' ">
        BACK HOME
      </button>
    </p>
  </center>
</body>
</html>