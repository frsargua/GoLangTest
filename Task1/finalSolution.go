package query

sqlQuert := "
SELECT 
    mt.name,
    mt.website_domain,
    t2.number_count
FROM 
    (
    SELECT 
        *,
        SUBSTRING(tb.website FROM '(?:.*://)?(?:www\.)?([^/?]*)')::VARCHAR AS website_domain
    FROM 
        "MY_TABLE" AS tb
    WHERE 
        tb.website IS NOT NULL AND 
        tb.website != ''
    ) AS mt
right outer JOIN (
    SELECT 
        SUBSTRING(t3.website FROM '(?:.*://)?(?:www\.)?([^/?]*)')::VARCHAR AS website_domain, 
        COUNT(*) as number_count
    FROM 
        "MY_TABLE" AS t3
    WHERE 
        t3.website IS NOT NULL AND 
        t3.website != ''
    GROUP BY 
        website_domain
    HAVING 
        COUNT(*) > 1
) AS t2
ON mt.website_domain = t2.website_domain

"