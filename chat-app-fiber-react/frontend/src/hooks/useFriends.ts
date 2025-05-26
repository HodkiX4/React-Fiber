import { useState } from "react";
import { useToastContext } from "../context/ToastContext";
import handleErrors from "../utils/errorUtil";
import {
  acceptFriendRequestApi,
  getPendingFriendRequestsApi,
  rejectFriendRequestApi,
  sendFriendRequestApi
} from "../api/friendsApi";
import { useFriendRequestStore } from "../store/friendReqStore";
import { useUsers } from "./useUsers";

export const useFriends = () => {
  const { getMe } = useUsers();
  const [sendloading, setSendLoading] = useState(false);
  const [acceptloading, setAcceptLoading] = useState(false);
  const [rejectloading, setRejectLoading] = useState(false);
  const [requestsloading, setRequestsLoading] = useState(false);
  const { showError, showSuccess } = useToastContext();
  const { setIncomingRequests, setOutgoingRequests, updateRequestStatus } = useFriendRequestStore();

  const sendFriendRequest = async (reciverId: string) => {
    setSendLoading(true);
    try {
      const data = await sendFriendRequestApi(reciverId);
      getPendingFriendRequests();
      showSuccess(data.message);
      
    } catch (error) {
      
      showError(handleErrors(error));
    } finally {
      setSendLoading(false);
    }
  }
  const acceptFriendRequest = async (friendReqId: string) => {
    setAcceptLoading(true);
    try {
      const data = await acceptFriendRequestApi(friendReqId);
      updateRequestStatus(friendReqId, "accepted");
      getPendingFriendRequests();
      getMe();
      showSuccess(data.message);
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setAcceptLoading(false);
    }
  }
  const rejectFriendRequest = async (friendReqId: string) => {
    setRejectLoading(true);
    try {
      const data = await rejectFriendRequestApi(friendReqId);
      updateRequestStatus(friendReqId, "rejected");
      getPendingFriendRequests();
      showSuccess(data.message);
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setRejectLoading(false);
    }
  }

  const getPendingFriendRequests = async () => {
    setRequestsLoading(true);
    try {
      const data = await getPendingFriendRequestsApi();
      setIncomingRequests(Array.isArray(data.incoming) ? data.incoming : []);
      setOutgoingRequests(Array.isArray(data.outgoing) ? data.outgoing : []);
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setRequestsLoading(false);
    }
  }

  return {
    sendloading, sendFriendRequest,
    acceptloading, acceptFriendRequest,
    rejectloading, rejectFriendRequest,
    requestsloading, getPendingFriendRequests
  }
}
