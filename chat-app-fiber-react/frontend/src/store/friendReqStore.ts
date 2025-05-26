import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { IFriendRequest } from "../models/friend.model";

interface FriendReqState {
  incomingRequests: IFriendRequest[];
  outgoingRequests: IFriendRequest[];
  setIncomingRequests: (requests: IFriendRequest[]) => void;
  setOutgoingRequests: (requests: IFriendRequest[]) => void;
  updateRequestStatus: (id: string, status: "accepted" | "rejected") => void;
}

export const useFriendRequestStore = create<FriendReqState>()(
  persist(
    (set) => ({
      incomingRequests: [],
      outgoingRequests: [],

      setIncomingRequests: (requests) => set({ incomingRequests: requests }),
      setOutgoingRequests: (requests) => set({ outgoingRequests: requests }),

      updateRequestStatus: (id, status) =>
        set((state) => ({
          incomingRequests: state.incomingRequests.map((r) =>
            r.id === id ? { ...r, status } : r
          ),
          outgoingRequests: state.outgoingRequests.map((r) =>
            r.id === id ? { ...r, status } : r
          ),
        })),
    }),
    {
      name: "friend-requests-storage",
      storage: createJSONStorage(() => localStorage),
    }
  )
);
