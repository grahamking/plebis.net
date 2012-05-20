<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head> 
	<meta http-equiv="content-type"	content="text/html; charset=utf-8" />
    <title>Plebis.Net. Say something.</title> 
	
	<link rel="stylesheet" type="text/css" href="http://yui.yahooapis.com/combo?2.6.0/build/reset-fonts-grids/reset-fonts-grids.css&2.6.0/build/base/base-min.css"> 

    <style type="text/css" media="screen">
		@import url( /style.css );
    </style>
	
</head>

<body class="yui-skin-sam">
	<div id="doc3" class="yui-t3"> 
		
		<div id="hd">
			<div class="alignleft">
				Plebis.Net
			</div>
			<div id="top_username">Speak your mind</div>
		</div>
		
		<div id="bd">
		
			<div id="sidebar" class="yui-b">

				<form id="message_form" name="message_form" method="post" action="/post.php">
					<fieldset>
						<input type="text" name="name" id="name" />
						<label for="name">Your name:</label>
					</fieldset>
					<fieldset>
						<label for="message">Your mind:</label>
						<textarea id="message" rows="15" cols="30" name="message"></textarea>
						<script type="text/javascript">
							document.write("<input type=hidden name='date' value='"+ new Date() +"' />");
						</script>
						<div class="action">
							<input type="submit" value="Write message" />
						</div>
					</fieldset>
				</form>

			</div>
		
			<div id="yui-main">
				<div class="yui-b">

					<div id="messages">
						<?php
						class Message {
							public $date, $who, $content;
						}
						
						$memcache = new Memcache;
						$memcache->connect('localhost', 11211);
						$msg_list_serialized = $memcache->get('plebis');
						if ( ! $msg_list_serialized ) {
							$filename = "/var/www/www.plebis.net/entries.dat";
							$fh = fopen($filename, "rt");
							$msg_list_serialized = fread($fh, filesize($filename));
							fclose($fh);
						}
						$msg_list = unserialize($msg_list_serialized);
						
						for ($i = count($msg_list) - 1; $i >= 0; $i--) {
						
							$msg = $msg_list[$i];
							echo '<div class="message">';
							echo '<h2>';
							echo $msg->who;
							echo '</h2>';
							echo nl2br($msg->content);
							echo '<div class="timestamp">'. $msg->when .'</div>';
							echo '</div>';
						}
						?>
					</div>
					
				</div>
			</div>

		</div>  
		
		<div id="ft">
			A <a href="http://www.darkcoding.net">darkcoding</a> production
		</div>
		
	</div> 
	
	<script type="text/javascript" src="http://yui.yahooapis.com/combo?2.6.0/build/yahoo-dom-event/yahoo-dom-event.js"></script>
	<script type="text/javascript">
		function init() {
			document.getElementById('name').focus();
		}
		YAHOO.util.Event.onDOMReady(init);
	</script>
</body>

</html>
