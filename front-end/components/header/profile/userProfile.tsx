import Image from "next/image";
import style from "@/styles/profile.module.css";

interface Props {
  url?: string;
  name?: string;
}

export default function UserProfileImage({ url, name = "?" }: Props) {
  const initials = name
    .trim()
    .split(" ")
    .map((part) => part[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();

  if (!url) {
    return (
      <div>
        <div
          className={style.imageCont}
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            background: "linear-gradient(135deg, var(--primary-theme), #1f1f1f)",
            color: "#fff",
            fontWeight: 700,
          }}
        >
          {initials}
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className={style.imageCont}>
        <Image
          className={style.userImage}
          src={url}
          fill
          priority
          sizes="(max-width: 768px) 100vw 50vw"
          alt="user profile image"
          loading="eager"
        />
      </div>
    </div>
  );
}