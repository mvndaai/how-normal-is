-- Core schema for How Normal Is (Supabase/Postgres)
create extension if not exists pgcrypto;

create table if not exists public.user_profiles (
  user_id uuid primary key references auth.users(id) on delete cascade,
  display_name text,
  can_comment boolean not null default false,
  can_moderate boolean not null default false,
  comment_banned_until timestamptz,
  comment_banned_forever boolean not null default false,
  blur_preference text not null default 'ask' check (blur_preference in ('ask','auto_unblur','skip_blur')),
  is_public boolean not null default false,
  public_share_questions boolean not null default false,
  public_share_answers boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists public.trait_groups (
  id uuid primary key default gen_random_uuid(),
  key text not null unique,
  label text not null
);

create table if not exists public.traits (
  id uuid primary key default gen_random_uuid(),
  name text not null unique,
  normalized_name text not null unique,
  group_id uuid references public.trait_groups(id),
  created_at timestamptz not null default now()
);

create table if not exists public.user_traits (
  user_id uuid not null references auth.users(id) on delete cascade,
  trait_id uuid not null references public.traits(id) on delete cascade,
  created_at timestamptz not null default now(),
  primary key(user_id, trait_id)
);

create table if not exists public.questions (
  id uuid primary key default gen_random_uuid(),
  asked_by uuid not null references auth.users(id) on delete cascade,
  prompt text not null,
  photo_path text,
  photo_blur_required boolean not null default false,
  status text not null default 'active' check (status in ('active','archived')),
  confusing_votes integer not null default 0,
  created_at timestamptz not null default now()
);

create table if not exists public.question_traits (
  question_id uuid not null references public.questions(id) on delete cascade,
  trait_id uuid not null references public.traits(id) on delete cascade,
  kind text not null default 'related' check (kind in ('related','disqualifier')),
  primary key(question_id, trait_id, kind)
);

create table if not exists public.question_answers (
  question_id uuid not null references public.questions(id) on delete cascade,
  user_id uuid not null references auth.users(id) on delete cascade,
  answer text not null check (answer in ('yes_for_me','sometimes_for_me','never_for_me')),
  because_trait_id uuid references public.traits(id),
  created_at timestamptz not null default now(),
  primary key(question_id, user_id)
);

create table if not exists public.question_confusing_votes (
  question_id uuid not null references public.questions(id) on delete cascade,
  user_id uuid not null references auth.users(id) on delete cascade,
  created_at timestamptz not null default now(),
  primary key(question_id, user_id)
);

create table if not exists public.comments (
  id uuid primary key default gen_random_uuid(),
  question_id uuid not null references public.questions(id) on delete cascade,
  user_id uuid not null references auth.users(id) on delete cascade,
  body text not null,
  sentiment_score double precision,
  status text not null default 'visible' check (status in ('visible','review','hidden')),
  created_at timestamptz not null default now()
);

create table if not exists public.comment_reviews (
  id uuid primary key default gen_random_uuid(),
  comment_id uuid not null references public.comments(id) on delete cascade,
  reviewer_id uuid not null references auth.users(id) on delete cascade,
  action text not null check (action in ('approve','reject','strong_reject')),
  strong_reject_reason text,
  created_at timestamptz not null default now()
);

create table if not exists public.moderator_accuracy_events (
  id uuid primary key default gen_random_uuid(),
  moderator_id uuid not null references auth.users(id) on delete cascade,
  approved_comment_id uuid references public.comments(id) on delete set null,
  later_flagged boolean not null default false,
  created_at timestamptz not null default now()
);
