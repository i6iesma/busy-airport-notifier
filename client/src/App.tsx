import { useEffect, useState } from 'react'
import './App.css'
const PUBLIC_KEY="BBCQIqZGOV52qqtx6NahpoKDj-H9gM6PxTeDvtqfwNCqxFUGqd3zYSQp2kM457CRoTk8XbjD_Y02kGeTh2I8DoE"

function App() {
  const [notificationButtonText, updateNotificationButtonText] = useState("Enable notifications")
  useEffect(() => {
    if ('serviceWorker' in navigator) {
      console.log("service worker in navigator")
     navigator.serviceWorker.register('/service-worker.js').then(registration => {
      registration.pushManager.permissionState().then((state) => console.log(state))
        }).catch(err => console.log(err))
      
    }
    }, [])

  const handleEnableNotifications = () => {
    const requestNotificationPermission = async () => {
      const permission = await Notification.requestPermission()
      if (permission == "granted") {
        console.log("permission granted")
        updateNotificationButtonText("Permission granted")
        subscribeUserToPush()
      }
      if (permission == "denied") {
        console.log("permission deined")
        updateNotificationButtonText("Permission denied, keep in mind notifications are needed to use this application")
      }
      if (permission == "default") {
        updateNotificationButtonText("Permission denied, keep in mind notifications are needed to use this application")
      }
    }
    requestNotificationPermission()
  }
  async function subscribeUserToPush() {
    const registration = await navigator.serviceWorker.ready
    try {
      const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlB64ToUint8Array(PUBLIC_KEY)
      })
      console.log(subscription)
      const p256Key = subscription.getKey("p256dh")
      const authKey = subscription.getKey("auth")
      if (p256Key && authKey) {
        const keyArray = Array.from(new Uint8Array(p256Key));
        const keyString = btoa(String.fromCharCode(...keyArray));
        console.log("p256dh key:", keyString);
        const authKeyArray = Array.from(new Uint8Array(authKey));
        const authKeyString = btoa(String.fromCharCode(...authKeyArray));
        console.log("auth key:", authKeyString)
      }
    } catch (err) {
        console.error(err)
    }


  }
  return (
    <>
    <button onClick={handleEnableNotifications}>{notificationButtonText}</button>
    </>
  
    
  )
}

export default App

 function urlB64ToUint8Array(base64String: string) {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding)
    .replace(/\-/g, "+")
    .replace(/_/g, "/");

  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}
