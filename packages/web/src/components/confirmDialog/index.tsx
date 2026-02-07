import React, { useRef, useState } from "react";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";

type ConfirmOptions = {
    title?: string;
    description?: string;
    confirmLabel?: string;
    cancelLabel?: string;
};

export function useConfirm() {
    const [open, setOpen] = useState(false);
    const resolverRef = useRef<((v: boolean) => void) | null>(null);
    const optsRef = useRef<ConfirmOptions>({});

    const confirm = (opts: ConfirmOptions = {}): Promise<boolean> => {
        return new Promise((resolve) => {
            optsRef.current = opts;
            resolverRef.current = resolve;
            setOpen(true);
        });
    };

    const handleConfirm = (v: boolean) => {
        setOpen(false);
        if (resolverRef.current) {
            resolverRef.current(v);
            resolverRef.current = null;
        }
    };

    const ConfirmDialog = (
        <AlertDialog
            open={open}
            onOpenChange={(o) => {
                if (!o) handleConfirm(false);
            }}
        >
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>{optsRef.current.title ?? "Confirm"}</AlertDialogTitle>
                    <AlertDialogDescription>
                        {optsRef.current.description ?? "Are you sure you want to continue?"}
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel onClick={() => handleConfirm(false)}>
                        {optsRef.current.cancelLabel ?? "Cancel"}
                    </AlertDialogCancel>
                    <AlertDialogAction onClick={() => handleConfirm(true)}>
                        {optsRef.current.confirmLabel ?? "Confirm"}
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );

    return { confirm, ConfirmDialog };
}