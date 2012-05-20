<?php

class Message {
	public $date, $who, $content;
}

$msg = new Message();
$msg->when = htmlspecialchars($_POST['date']);
$msg->who = htmlspecialchars($_POST['name']);
$msg->content = htmlspecialchars($_POST['message']);

if ( empty($msg->when) || strlen($msg->when) < 10 ) {
    echo '<h1>Sorry, we think you might be a spammer</h1>';
    return;
}
    
$memcache = new Memcache();
$memcache->connect('localhost', 11211);

$current_serialized = $memcache->get('plebis');
if ( ! $current_serialized ) {
	$filename = "/var/www/www.plebis.net/entries.dat";
	$fh = fopen($filename, "rt");
	$current_serialized = fread($fh, filesize($filename));
	fclose($fh);
}
$current = unserialize($current_serialized);
if ( empty($current) ) {
	$current = array();
}

$current[] = $msg;

if ( count($current) > 50 ) {
	$current = array_slice($current, count($current)-50 , 50);
}

$current_str = serialize($current);

$memcache->set('plebis', $current_str, false, 0);

$fh = fopen("/var/www/www.plebis.net/entries.dat", "wt");
fwrite($fh, $current_str);
fclose($fh);

header('Location: /');

?>
