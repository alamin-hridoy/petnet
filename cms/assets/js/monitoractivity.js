function monitorActivity() {
    
  // Set the javascript window events to monitor activity
  window.onload = activityDetected;      // detects initial load / reload
  window.onkeypress = activityDetected;  // detects key presses
  window.onclick = activityDetected;     // detects mouse clicks
  window.onmousemove = activityDetected; // detects mouse movement
  window.onscroll = activityDetected;    // detects arrow key scrolling
  window.onmousedown = activityDetected; // detects touchscreen interaction
  
  function onNoActivity() {
    window.location.href = '/logout';
  }

  var timer;
  var timout = 900000; // milliseconds
  function activityDetected() {
          clearTimeout(timer);
          timer = setTimeout(onNoActivity, timout);
      }
  }
  monitorActivity();