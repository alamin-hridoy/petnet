// off canvas menu
function openNav() {
  document.getElementById("mobile-nav").classList.add("opened");
}
function closeNav() {
  document.getElementById("mobile-nav").classList.remove("opened");
}

let _urls
async function getUrls() {
  if (_urls != null) {
    return _urls;
  }

  const response = await fetch("/urls.json")
  if (response.status !== 200) {
    throw "no URLs";
  }
  const urls = await response.json();
  _urls = urls
  return _urls
}


// Toggle show password
//
function showPassword() {
  var password = document.querySelector("input.password");
  password.type = password.type === "password" ? "text" : "password";
}

// Password Meter
//
var Password = document.getElementById("new-password");

var strengthbar = document.getElementById("password-meter");
var display = document.getElementsByClassName("password-textbox")[0];

if (Password) {
  Password.addEventListener("keyup", function () {
    checkpassword(Password.value);
  });
}

function preventSpace(event) {
  if (event.code === "Space") {
    event.preventDefault();
  }
}

function checkpassword(password) {
  var strength = 0;
  if (password.match(/[a-z]+/)) {
    strength += 1;
    display.innerHTML = "Better to add at least 1 uppercase case letter [A-Z]";
  }
  if (password.match(/[A-Z]+/)) {
    strength += 1;
    display.innerHTML = "Better to add at least 1 number [1, 2].";
  }
  if (password.match(/[0-9]+/)) {
    strength += 1;
    display.innerHTML = "Better to add at least 1 spacial letter [!, #].";
  }
  if (password.match(/[$@#&!]+/)) {
    strength += 1;
  }

  if (password.length < 8) {
    display.innerHTML =
      "Password must contain between 8 and 16 characters in length.";
  }

  if (password.length > 16) {
    display.innerHTML = "maximum number of characters is 16";
  }

  switch (strength) {
    case 0:
      strengthbar.setAttribute("data", "0");
      break;

    case 1:
      strengthbar.setAttribute("data", "25");
      break;

    case 2:
      strengthbar.setAttribute("data", "50");
      break;

    case 3:
      strengthbar.setAttribute("data", "75");
      break;

    case 4:
      strengthbar.setAttribute("data", "100");
      break;
  }
}

// Home Tabs
//
// tabs
(function () {
  var tabs = document.querySelectorAll(".tabs");

  [].forEach.call(tabs, function (tab) {
    tabNav(tab);
  });
})();
function tabNav(tab) {

  var tabLinks = tab.querySelectorAll(".tablinks");
  var tabContent = tab.querySelectorAll(".tabcontent");

  tabLinks.forEach(function (el) {
    el.addEventListener("click", openTabs);
  });

  function openTabs(el) {
    var btnTarget = el.currentTarget;
    var tab = btnTarget.dataset.tab;

    tabContent.forEach(function (el) {
      el.classList.remove("active");
    });

    tabLinks.forEach(function (el) {
      el.classList.remove("active");
    });

    document.querySelector("#" + tab).classList.add("active");
    btnTarget.classList.add("active");
  }
}

// Tabs Content Carousel
//
(function () {
  var carousels = document.querySelectorAll(".carousel");

  [].forEach.call(carousels, function (carousel) {
    carouselize(carousel);
  });
})();

function carouselize(carousel) {
  var slides = carousel.querySelector(".slides");
  var slideListWidth = 0;
  var slideListSteps = 0;
  var allSlides = carousel.querySelectorAll(".slide");
  var slideAmount = 0;
  var slideAmountVisible = carousel.getAttribute("data-slide");
  var carouselPrev = carousel.querySelector(".carousel-prev");
  var carouselNext = carousel.querySelector(".carousel-next");

  var carouselContainer = document.querySelector(".carousel-container").offsetWidth - 32;
  var slideWidth = carouselContainer / slideAmountVisible;

  var tabletScreen = window.matchMedia("(max-width: 768px)");
  var mobileScreen = window.matchMedia("(max-width: 375px)");

  var slideAmountMobile = carousel.getAttribute("data-mobile-slides");
  var mobileslideWidth = carouselContainer / slideAmountMobile;
  var slideAmountTablet = carousel.getAttribute("data-tablet-slides");
  var tabletslideWidth = carouselContainer / slideAmountTablet;

  var allDots = carousel.querySelectorAll(".dots span");

  //Count all the slides
  [].forEach.call(allSlides, function (slide) {
    slideAmount++;

    if (mobileScreen.matches) {
      slideListWidth += mobileslideWidth + 16;
      slides.style.width = slideListWidth + "px";
     }
    else if (tabletScreen.matches) {
      slideListWidth += tabletslideWidth + 16;
      slides.style.width = slideListWidth + "px";
     }
     else {
      slideListWidth += slideWidth + 16;
      slides.style.width = slideListWidth + "px";
     }

    carousel.querySelectorAll(".slide").forEach((slide) => {
      if (mobileScreen.matches) {
        slide.style.width = mobileslideWidth + "px";
       }
       else if (tabletScreen.matches) {
        slide.style.width = tabletslideWidth + "px";
       }
       else {
        slide.style.width = slideWidth + "px";
       }
    });
  });

  // Dots navigation
  [].forEach.call(allDots, function (dot) {
    var dataDot = dot.getAttribute("data-dot");

    dot.onclick = function () {
      slides.style.transform =
        "translateX(-" + slideWidth * (dataDot * slideAmountVisible) + "px)";
      removeDotsActive();
      dot.classList.add("active");
    };
  });

  carouselNext.onclick = function () {
    if ((mobileScreen.matches) && (slideListSteps < slideAmount - slideAmountMobile)) {
      slideListSteps++;
      moveSlideList();
     }
    else if ((tabletScreen.matches) && (slideListSteps < slideAmount - slideAmountTablet)) {
      slideListSteps++;
      moveSlideList();
     }
     else if (slideListSteps < slideAmount - slideAmountVisible) {
      slideListSteps++;
      moveSlideList();
    }
  };
  carouselPrev.onclick = function () {
    if (slideListSteps > 0) {
      slideListSteps--;
      moveSlideList();
    }
  };

  // Move the carousels product-list
  function moveSlideList() {
    if (mobileScreen.matches) {
      slides.style.transform = "translateX(-" + mobileslideWidth * slideListSteps + "px)";
     }
    else if (tabletScreen.matches) {
      slides.style.transform = "translateX(-" + tabletslideWidth * slideListSteps + "px)";
     }
     else {
      slides.style.transform = "translateX(-" + slideWidth * slideListSteps + "px)";
     }
  }

  function removeDotsActive() {
    for (i = 0; i < allDots.length; i++) {
      allDots[i].classList.remove("active");
    }
  }
}

// Accordion
//
const accs = document.getElementsByClassName("accordion");

function toggleAccordion(i) {
  let acc = accs[i];
  acc.classList.toggle("active");

  let panel = acc.nextElementSibling;
  if (panel.style.maxHeight) {
    panel.style.maxHeight = null;
  } else {
    panel.style.maxHeight = panel.scrollHeight + "px";
  }
}

for (let i = 0; i < accs.length; i++) {
  let acc = accs[i]
  acc.addEventListener("click", function () {
    toggleAccordion(i)
  });

  if (acc.classList.contains("opened-on-start")) {
    toggleAccordion(i)
  }
}

function displayError(msg) {
  document.getElementById("login-error").classList.remove("hidden");
  document.getElementById("error-msg").innerText = msg;
}

async function signup(event) {
  event.preventDefault();

  const email = document.getElementById("email").value;
  const password = document.getElementById("new-password").value;

  const urls = getUrls()
  const body = { name: "N/A", email, password };
  const signUpUrl = urls.user_mgm + "/v1/signup";
  const response = await fetch(signUpUrl, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
  if (response.status === 200) {
    window.location.replace("/");
  } else {
    const content = await response.json();
    displayError(content.message);
  }
}

function clearErrorMessages(errorMsgs) {
  if (event) {
    event.target.classList.remove('border-red')
  }
  errorMsgs.forEach(id => {
    errorMsgEl = document.getElementById(id);
    if (errorMsgEl) {
      errorMsgEl.classList.add('hidden');
    }
  })
}
