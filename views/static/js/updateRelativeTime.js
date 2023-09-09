    // Function to calculate and update relative time
    function updateRelativeTime() {
      const timeElements = document.querySelectorAll(".time-ago");
      timeElements.forEach((element) => {
        const formattedTime = element.getAttribute("data-timestamp");
        const trimmedTime = formattedTime.split(".")[0];
        const timestamp = Math.floor(new Date(trimmedTime).getTime() / 1000);
        const currentTime = Math.floor(Date.now() / 1000);
        const timeDifference = currentTime - timestamp;
        
        let timeString;
        if (timeDifference < 60) {
          timeString = `${timeDifference}s ago`;
        } else if (timeDifference < 3600) {
          const minutes = Math.floor(timeDifference / 60);
          timeString = `${minutes}min ago`;
        } else if (timeDifference < 86400) {
          const hours = Math.floor(timeDifference / 3600);
          timeString = `${hours}h ago`;
        } else if (timeDifference < 172800) {
          timeString = "1day ago";
        } else {
          const days = Math.floor(timeDifference / 86400);
          timeString = `${days}days ago`;
        }

        element.textContent = timeString;
      });
    }

    // Call the function initially
    updateRelativeTime();

    // Call the function every minute to keep updating the time difference
    setInterval(updateRelativeTime, 60000); // 60000 milliseconds = 1 minute