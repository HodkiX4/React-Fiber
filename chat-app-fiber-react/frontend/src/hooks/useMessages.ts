import { useState } from "react";
import { useRoomsStore } from "../store/roomStore";
import { deleteMessagesApi, getMessagesApi, sendMessageApi } from "../api/messageApi";
import handleErrors from "../utils/errorUtil";
import { useToastContext } from "../context/ToastContext";

export const useMessages = () => {
  const {
    activeRoomId,
    addMessage,
    removeMessage,
    setMessages
  } = useRoomsStore();
  const { showSuccess, showError } = useToastContext();
  const [sendLoading, setSendLoading] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [getLoading, setGetLoading] = useState(false);
  
  const getMessages = async () => {
    if(getLoading) return;
    if(activeRoomId == null) {
      showError("No room is selected");
      return;
    }
    setGetLoading(true);
    try {
      const data = await getMessagesApi(activeRoomId);
      setMessages(activeRoomId, data.messages);
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setGetLoading(false);
    }
  }

  const sendMessage = async (content: string) => {
    if(sendLoading) return;
    if(activeRoomId == null) {
      showError("No room is selected");
      return;
    } 

    setSendLoading(true);
    try {
      const data = await sendMessageApi({ roomId: activeRoomId, content });
      addMessage(activeRoomId, data.message);
      showSuccess("Message sent!");
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setSendLoading(false);
    }
  }
  
  const deleteMessage = async (messageId: string) => {
    if(deleteLoading) return;
    if(activeRoomId == null) {
      showError("No room is selected");
      return;
    }
    setDeleteLoading(true)
    try {
      const data = await deleteMessagesApi(messageId);
      removeMessage(activeRoomId, messageId)
      showSuccess(data.message);
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setDeleteLoading(false);
    }
  }


  return {
    getMessages, getLoading,
    sendMessage, sendLoading,
    deleteMessage, deleteLoading
  }
}