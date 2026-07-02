import DefaultLayout from "@/components/layouts/defaultLayout";
import NotificationsPage from "./notifications";

export default function NotificationUI() {
    console.log("NotificationUI rendered");

    return (
        <DefaultLayout>
            <NotificationsPage />
        </DefaultLayout>
    );
}