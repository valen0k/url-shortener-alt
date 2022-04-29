CREATE TABLE IF NOT EXISTS public.shortener (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    hash_url VARCHAR(10) NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
